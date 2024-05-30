package main

import (
	"context"
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

	err = run.Get(context.Background(), nil)
	if err != nil {
		return err
	}

	return nil
}

func FifoWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)
	messages := []string{}
	waiting := true

	timerFuture := workflow.NewTimer(ctx, 5*time.Second)
	selector.AddFuture(timerFuture, func(f workflow.Future) {
		waiting = false
	})

	signalChn := workflow.GetSignalChannel(ctx, signalName)
	selector.AddReceive(signalChn, func(c workflow.ReceiveChannel, more bool) {
		var msg string
		c.Receive(ctx, &msg)
		messages = append(messages, msg)
	})

	for waiting || selector.HasPending() {
		selector.Select(ctx)
	}

	for _, msg := range messages {
		logger.Info("Received:", "msg", msg)
		// execute Activities
	}

	return "done", nil
}
