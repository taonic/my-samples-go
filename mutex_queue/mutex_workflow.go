package mutex_queue

import (
	"fmt"
	"slices"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	// AcquireLockSignalName signal channel name for lock acquisition
	AcquireLockSignalName = "acquire-lock-event"
	// RequestLockSignalName channel name for request lock
	RequestLockSignalName = "request-lock-event"
)

func MutexWorkflowWithCancellation(
	ctx workflow.Context,
	namespace string,
	resourceID string,
	unlockTimeout time.Duration,
) error {
	logger := workflow.GetLogger(ctx)
	currentWorkflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	if currentWorkflowID == "default-test-workflow-id" {
		// unit testing hack, see https://github.com/uber-go/cadence-client/issues/663
		_ = workflow.Sleep(ctx, 10*time.Millisecond)
	}
	logger.Info("started", "currentWorkflowID", currentWorkflowID)

	requestLockCh := workflow.GetSignalChannel(ctx, RequestLockSignalName)
	completeCh := workflow.NewChannel(ctx)
	done := false
	queue := []string{}

	for {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(completeCh, func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, &done)
		})
		selector.AddReceive(requestLockCh, func(c workflow.ReceiveChannel, more bool) {
			var senderID string
			c.Receive(ctx, &senderID)
			queue = append(queue, senderID)
			// Start a goroutine per sender
			workflow.Go(ctx, func(ctx workflow.Context) {
				thisSenderID := senderID
				var index int
				workflow.Await(ctx, func() bool {
					index = slices.Index(queue, thisSenderID)
					isLast := index == len(queue)-1
					// Block the last element in the queue unless it's the only one left.
					return !isLast || len(queue) == 1
				})
				switch index {
				case -1:
					return
				case 0:
					unblockSender(ctx, senderID, unlockTimeout)
				default:
					cancelSender(ctx, senderID)
				}
				// Remove the unblocked or cancelled sender
				queue = append(queue[:index], queue[index+1:]...)
				// Try to complete if queue is empty
				if len(queue) == 0 {
					completeCh.Send(ctx, true)
				}
			})
		})
		if done && requestLockCh.Len() == 0 {
			return nil
		}
		selector.Select(ctx)
	}
}

func cancelSender(ctx workflow.Context, senderWorkflowID string) {
	logger := workflow.GetLogger(ctx)
	err := workflow.RequestCancelExternalWorkflow(ctx, senderWorkflowID, "").Get(ctx, nil)
	if err != nil {
		logger.Info("CancelExternalWorkflow error", "Error", err)
	}
}

func unblockSender(ctx workflow.Context, senderWorkflowID string, unlockTimeout time.Duration) {
	logger := workflow.GetLogger(ctx)
	var releaseLockChannelName string
	_ = workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return generateUnlockChannelName(senderWorkflowID)
	}).Get(&releaseLockChannelName)
	logger.Info("generated release lock channel name", "releaseLockChannelName", releaseLockChannelName)
	// Send release lock channel name back to a senderWorkflowID, so that it can
	// release the lock using release lock channel name
	err := workflow.SignalExternalWorkflow(ctx, senderWorkflowID, "", AcquireLockSignalName, releaseLockChannelName).Get(ctx, nil)
	if err != nil {
		// .Get(ctx, nil) blocks until the signal is sent.
		// If the senderWorkflowID is closed (terminated/canceled/timeouted/completed/etc), this would return error.
		// In this case we release the lock immediately instead of failing the mutex workflow.
		// Mutex workflow failing would lead to all workflows that have sent requestLock will be waiting.
		logger.Info("SignalExternalWorkflow error", "Error", err)
	}
	logger.Info("signaled external workflow")
	var ack string
	workflow.GetSignalChannel(ctx, releaseLockChannelName).ReceiveWithTimeout(ctx, unlockTimeout, &ack)
	logger.Info("release signal received: " + ack)
}

// generateUnlockChannelName generates release lock channel name
func generateUnlockChannelName(senderWorkflowID string) string {
	return fmt.Sprintf("unlock-event-%s", senderWorkflowID)
}
