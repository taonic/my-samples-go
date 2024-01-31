package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	samples "github.com/taonic/my-samples-go"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientOptions, err := samples.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}

	c, err := client.Dial(clientOptions)
	if err != nil {
		return err
	}
	defer c.Close()

	// Start worker
	log.Print("Running worker with limited concurrent activities")
	const taskQueue = "my-task-queue"
	w := worker.New(c, taskQueue, worker.Options{
		// Only allow so many at once
		MaxConcurrentActivityExecutionSize: 2,
	})
	w.RegisterWorkflow(MyWorkflow)
	w.RegisterActivity(MyActivity)
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	// Run workflow
	log.Print("Running workflow")
	run, err := c.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{TaskQueue: taskQueue},
		MyWorkflow,
	)
	if err != nil {
		return err
	} else if err := run.Get(context.Background(), nil); err != nil {
		return err
	}
	log.Printf("Workflow done")
	return nil
}

func MyActivity(ctx context.Context, name string) error {
	sleepFor := time.Duration(0) * time.Second
	log.Printf("Activity %v started, sleeping for %v", name, sleepFor)
	defer log.Printf("Activity %v stopped", name)
	time.Sleep(sleepFor)
	return activity.ErrResultPending
}

func MyWorkflow(ctx workflow.Context) error {
	// Start 20 activities
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: 100 * time.Second})
	futs := make([]workflow.Future, 20)
	for i := range futs {
		futs[i] = workflow.ExecuteActivity(ctx, MyActivity, fmt.Sprintf("activity%v", i+1))
	}
	// Wait for all 20 to be done
	for i, fut := range futs {
		activityName := fmt.Sprintf("activity%v", i+1)
		if err := fut.Get(ctx, nil); err != nil {
			workflow.GetLogger(ctx).Error(activityName+" failed", "error", err)
		} else {
			workflow.GetLogger(ctx).Info(activityName + " complete")
		}
	}
	return nil
}
