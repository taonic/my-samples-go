package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
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

	c, err := client.Dial(client.Options{})
	if err != nil {
		return err
	}
	defer c.Close()

	// Start worker
	const taskQueue = "my-task-queue-1"
	w := worker.New(c, taskQueue, worker.Options{})

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
		"hello",
	)
	if err != nil {
		return err
	} else if err := run.Get(context.Background(), nil); err != nil {
		return err
	}
	log.Printf("Workflow done")
	return nil
}

func generateUUID(ctx workflow.Context) (string, error) {
	var generatedUUID string

	err := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return uuid.NewString()
	}).Get(&generatedUUID)
	if err != nil {
		return "", err
	}

	return generatedUUID, nil
}

func MyActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	return "Hello " + name + "!", nil
}

// Workflow is a Hello World workflow definition.
func MyWorkflow(ctx workflow.Context, name string) (string, error) {
	uid := ""
	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	v := workflow.GetVersion(ctx, "hello-version", workflow.DefaultVersion, 1)
	if v == 1 {
		var err error
		uid, err = generateUUID(ctx)
		if err != nil {
			logger.Error("failed to generated uuid", "Error", err)
			return "", err
		}

		logger.Info("generated uuid", "uuid-val", uid)
	}

	ao := workflow.ActivityOptions{StartToCloseTimeout: 1 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, MyActivity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)

	return uid, nil
}
