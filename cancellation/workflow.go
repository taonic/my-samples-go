package cancellation

import (
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART samples-go-cancellation-workflow-definition
// YourWorkflow is a Workflow Definition that shows how it can be canceled.
func YourWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		HeartbeatTimeout:    5 * time.Second,
		WaitForCancellation: true,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("cancel workflow started")
	var a *Activities // Used to call Activities by function pointer
	defer func() {
		if !errors.Is(ctx.Err(), workflow.ErrCanceled) {
			return
		}

		// When the Workflow is canceled, it has to get a new disconnected context to execute any Activities
		newCtx, _ := workflow.NewDisconnectedContext(ctx)
		err := workflow.ExecuteActivity(newCtx, a.CleanupActivity).Get(ctx, nil)
		if err != nil {
			logger.Error("CleanupActivity failed", "Error", err)
		}
	}()

	completeCh := workflow.NewChannel(ctx)

	workflow.Go(ctx, func(ctx workflow.Context) {
		var result string
		err := workflow.ExecuteActivity(ctx, a.ActivityToBeCanceled).Get(ctx, &result)
		logger.Info(fmt.Sprintf("ActivityToBeCanceled returns %v, %v", result, err))

		err = workflow.ExecuteActivity(ctx, a.ActivityToBeSkipped).Get(ctx, nil)
		logger.Error("Error from ActivityToBeSkipped", "Error", err)

		completeCh.Send(ctx, true)
	})

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(completeCh, func(c workflow.ReceiveChannel, more bool) {
		logger.Info("Workflow execution complete.")
	})
	selector.AddReceive(ctx.Done(), func(c workflow.ReceiveChannel, more bool) {
		logger.Info("Workflow cancelled.")
	})
	selector.Select(ctx)

	return nil
}
