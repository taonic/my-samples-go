package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var (
	c          client.Client
	workflowID = uuid.New().String()
	taskQueue  = "my-task-queue" + uuid.New().String()
)

const (
	signalName = "my-signal"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error

	c, err = client.Dial(client.Options{})
	if err != nil {
		return err
	}
	defer c.Close()

	// Start worker
	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(FifoWorkflow)
	w.RegisterActivity(Activity)
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	// Run workflows
	wo := client.StartWorkflowOptions{
		TaskQueue: taskQueue,
	}
	var run client.WorkflowRun

	concurrency := 10
	maxMessages := 100
	// Send concurrent signals in batches
	for j := 0; j < 20; j++ {
		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			myCounter := j*10 + i
			go func() {
				defer wg.Done()
				run, err = c.SignalWithStartWorkflow(
					ctx,
					workflowID,
					signalName,
					fmt.Sprintf("belated message %d", myCounter),
					wo,
					FifoWorkflow,
					concurrency,
					maxMessages,
					[]string{},
				)
				if err != nil {
					panic(err)
				}
			}()
		}
		wg.Wait()
		time.Sleep(1 * time.Second)
	}

	err = run.Get(context.Background(), nil)
	if err != nil {
		return err
	}

	return nil
}

type Listener struct {
	Messages []string
	Max      int
	Done     bool
}

func (l *Listener) Listen(ctx workflow.Context) {
	chn := workflow.GetSignalChannel(ctx, signalName)
	processed := 0
	selector := workflow.NewSelector(ctx)
	selector.AddReceive(chn, func(c workflow.ReceiveChannel, more bool) {
		var msg string
		c.Receive(ctx, &msg)
		processed++
		l.Messages = append(l.Messages, msg)
	})

	for {
		if processed > l.Max && chn.Len() == 0 {
			l.Done = true
			return
		}

		selector.Select(ctx)
	}
}

func Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Hello", "name", name)
	return "Hello " + name + "!", nil
}

func FifoWorkflow(ctx workflow.Context, concurrency int, maxMessages int, messages []string) (string, error) {
	logger := workflow.GetLogger(ctx)
	l := Listener{
		Messages: messages,
		Max:      maxMessages,
	}
	workflow.Go(ctx, l.Listen)

	for {
		// Wait for listener to populate messages
		err := workflow.Await(ctx, func() bool {
			return len(l.Messages) > 0
		})
		if err != nil {
			return "", err
		}

		// Execute activities concurrently
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: time.Second * 5,
		})
		futures := []workflow.Future{}
		for len(l.Messages) > 0 && len(futures) < concurrency {
			var msg string
			msg, l.Messages = l.Messages[0], l.Messages[1:]
			f := workflow.ExecuteActivity(ctx, Activity, msg)
			futures = append(futures, f)
		}

		// Wait for the concurrent batch to complete
		for _, f := range futures {
			var result string
			err := f.Get(ctx, &result)
			if err != nil {
				return "", err
			}
			logger.Info("Result", result)
		}

		// Continue as new if needed
		if l.Done {
			return "", workflow.NewContinueAsNewError(ctx, FifoWorkflow, concurrency, maxMessages, l.Messages)
		}
	}
}
