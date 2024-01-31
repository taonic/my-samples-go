package main

import (
	"context"
	"log"
	"time"

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

	c, err := client.Dial(client.Options{})
	if err != nil {
		return err
	}
	defer c.Close()

	// Start worker
	const taskQueue = "my-task-queue"
	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(MyWorkflow)
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

func MyWorkflow(ctx workflow.Context) error {
	for i := 0; i < 12; i++ {
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
