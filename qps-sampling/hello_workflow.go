package qpssampling

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// HelloWorldActivity simple activity that returns a greeting
func HelloWorldActivity(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello %s!", name), nil
}

// HelloWorldWorkflow executes a configurable number of hello world activities
func HelloWorldWorkflow(ctx workflow.Context, numActivities int) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for i := 0; i < numActivities; i++ {
		var result string
		err := workflow.ExecuteActivity(ctx, HelloWorldActivity, fmt.Sprintf("World-%d", i+1)).Get(ctx, &result)
		if err != nil {
			return err
		}
		workflow.GetLogger(ctx).Info("Activity completed", "result", result)
	}
	return nil
}
