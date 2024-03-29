package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/taonic/my-samples-go/lib"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	clientOptions, err := lib.ParseClientOptionFlags(os.Args[1:])
	clientOptions.MetricsHandler = sdktally.NewMetricsHandler(lib.NewPrometheusScope(prometheus.Configuration{
		ListenAddress: "0.0.0.0:9091",
		TimerType:     "histogram",
	}))

	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := client.Dial(clientOptions)
	if err != nil {
		return err
	}
	defer c.Close()

	// Start worker
	var taskQueue = "my-task-queue" + uuid.New().String()
	w := worker.New(c, taskQueue, worker.Options{
		MaxConcurrentWorkflowTaskPollers: 80, // bumped from the default 2 to support the highly concurrent workload
		MaxConcurrentActivityTaskPollers: 80,
	})
	w.RegisterWorkflow(MyGrandParentWorkflow)
	w.RegisterWorkflow(MyParentWorkflow)
	w.RegisterWorkflow(MyChildWorkflow)
	w.RegisterActivity(MyActivity)
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	// Run workflow
	log.Print("Running workflow")
	combinations := [][]int{
		// partitions, size, activities per WF
		{1, 200, 1},
		{2, 100, 1},
		{4, 50, 1},
		{8, 25, 1},
		{20, 10, 1},
		{40, 5, 1},
	}
	results := []string{}
	for _, combination := range combinations {
		partitions, partitionSize, activities := combination[0], combination[1], combination[2]
		then := time.Now()
		wfID := fmt.Sprintf("benchmark/partitions:%d/partitionSize:%d/activities:%d", partitions, partitionSize, activities)
		run, err := c.ExecuteWorkflow(
			ctx,
			client.StartWorkflowOptions{ID: wfID, TaskQueue: taskQueue},
			MyGrandParentWorkflow,
			partitions,
			partitionSize,
			activities,
		)
		if err != nil {
			return err
		} else if err := run.Get(context.Background(), nil); err != nil {
			return err
		}
		results = append(results, fmt.Sprintf("%s was done in %s", wfID, time.Since(then)))
	}

	for _, result := range results {
		fmt.Println(result)
	}
	return nil
}

func MyGrandParentWorkflow(ctx workflow.Context, partitions, partitionSize, activities int) (string, error) {
	logger := workflow.GetLogger(ctx)

	var futures []workflow.ChildWorkflowFuture
	for i := 0; i < partitions; i++ {
		cwo := workflow.ChildWorkflowOptions{
			WorkflowID: "ParentWorkflow" + uuid.New().String(),
		}
		ctx = workflow.WithChildOptions(ctx, cwo)
		future := workflow.ExecuteChildWorkflow(ctx, MyParentWorkflow, partitionSize, activities)
		futures = append(futures, future)
	}

	var result string
	for _, future := range futures {
		err := future.Get(ctx, &result)
		if err != nil {
			return "", err
		}
		logger.Debug("Child execution completed.", "Result", result)
	}

	return "", nil
}

func MyParentWorkflow(ctx workflow.Context, numWorkflows, numActivities int) (string, error) {
	logger := workflow.GetLogger(ctx)
	var futures []workflow.ChildWorkflowFuture

	for i := 0; i < numWorkflows; i++ {
		cwo := workflow.ChildWorkflowOptions{
			WorkflowID: "ChildWorkflow" + uuid.New().String(),
		}
		ctx = workflow.WithChildOptions(ctx, cwo)
		future := workflow.ExecuteChildWorkflow(ctx, MyChildWorkflow, numActivities)
		futures = append(futures, future)
	}

	var result string
	for _, future := range futures {
		err := future.Get(ctx, &result)
		if err != nil {
			return "", err
		}
		logger.Debug("Child execution completed.", "Result", result)
	}

	return "", nil
}

func MyChildWorkflow(ctx workflow.Context, numActivities int) (string, error) {
	logger := workflow.GetLogger(ctx)

	var futures []workflow.Future

	for i := 0; i < numActivities; i++ {
		ao := workflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
		}
		ctx = workflow.WithActivityOptions(ctx, ao)
		future := workflow.ExecuteActivity(ctx, MyActivity, i)
		futures = append(futures, future)
	}

	var result string
	for _, future := range futures {
		err := future.Get(ctx, &result)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return "", err
		}
	}

	return "", nil
}

func MyActivity(ctx context.Context, i int) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Debug("Activity", i)
	return fmt.Sprintf("Hello %i!", i), nil
}
