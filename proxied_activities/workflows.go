package proxied_activities

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func HelloProxyWorkflow(ctx workflow.Context, message string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	var a ProxiedActivity
	err := workflow.ExecuteActivity(ctx, a.ProxyActivity, message).Get(ctx, &result)
	return result, err
}
