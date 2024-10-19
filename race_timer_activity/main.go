package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var (
	c          client.Client
	wg         *sync.WaitGroup
	workflowID = uuid.New().String()
	taskQueue  = "my-task-queue" + uuid.New().String()
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
	w.RegisterWorkflow(RaceTimerActivityWorkflow)
	w.RegisterActivity(SampleActivity)
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	// Start WF
	log.Print("Running workflow")
	run, err := c.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        workflowID,
			TaskQueue: taskQueue,
		},
		RaceTimerActivityWorkflow,
	)
	if err != nil {
		return err
	}

	err = run.Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func RaceTimerActivityWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)

	ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, ao)
	selector.AddFuture(workflow.ExecuteActivity(ctx, SampleActivity), func(f workflow.Future) {
		var output string
		if err := f.Get(ctx, &output); err == nil {
			logger.Info(fmt.Sprintf("-- Activity completed: %s --", output))
		}
	})
	selector.AddFuture(workflow.NewTimer(ctx, 1*time.Second), func(f workflow.Future) {
		logger.Info("-- Timer fired --")
	})

	for i := 0; i < 2; i++ {
		selector.Select(ctx)
	}

	return nil
}

func SampleActivity(ctx context.Context) (string, error) {
	time.Sleep(2 * time.Second)
	return "done", nil
}
