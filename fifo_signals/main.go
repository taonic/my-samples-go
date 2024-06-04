package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
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
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	// Run workflows
	wo := client.StartWorkflowOptions{
		TaskQueue: taskQueue,
	}
	var run client.WorkflowRun

	// Send first batch
	for i := 0; i < 5; i++ {
		run, err = c.SignalWithStartWorkflow(
			ctx,
			workflowID,
			signalName,
			strconv.Itoa(i),
			wo,
			FifoWorkflow,
		)
		if err != nil {
			return err
		}
	}

	// Wait for
	time.Sleep(2 * time.Second)

	// Send second batch
	for i := 0; i < 5; i++ {
		run, err = c.SignalWithStartWorkflow(
			ctx,
			workflowID,
			signalName,
			fmt.Sprintf("belated message %d", i),
			wo,
			FifoWorkflow,
		)
		if err != nil {
			return err
		}
	}

	err = run.Get(context.Background(), nil)
	if err != nil {
		return err
	}

	return nil
}

type Listener struct {
	Messages []string
}

func (l *Listener) Listen(ctx workflow.Context) {
	selector := workflow.NewSelector(ctx)
	signalChn := workflow.GetSignalChannel(ctx, signalName)
	selector.AddReceive(signalChn, func(c workflow.ReceiveChannel, more bool) {
		var msg string
		c.Receive(ctx, &msg)
		l.Messages = append(l.Messages, msg)
	})

	for {
		selector.Select(ctx)
	}
}

func FifoWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)
	l := Listener{Messages: []string{}}
	workflow.Go(ctx, l.Listen)

	for {
		// Wait for listener to populate messages
		err := workflow.Await(ctx, func() bool {
			return len(l.Messages) > 0
		})
		if err != nil {
			return "", err
		}

		// Process messages
		for len(l.Messages) > 0 {
			var msg string
			msg, l.Messages = l.Messages[0], l.Messages[1:]
			logger.Info("Received", "msg", msg)
			// execute your Activities
		}

		// TODO add continueAsNew
	}
}
