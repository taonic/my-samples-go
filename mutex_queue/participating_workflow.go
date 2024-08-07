package mutex_queue

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	UnlockFunc func() error

	Mutex struct {
		currentWorkflowID string
		lockNamespace     string
	}

	ContextKey string
)

const (
	TaskQueue                   = "mutex_queue"
	ClientContextKey ContextKey = "Client"
)

// NewMutex initializes mutex
func NewMutex(currentWorkflowID string, lockNamespace string) *Mutex {
	return &Mutex{
		currentWorkflowID: currentWorkflowID,
		lockNamespace:     lockNamespace,
	}
}

func (s *Mutex) LockWithCancellation(ctx workflow.Context,
	resourceID string, unlockTimeout time.Duration) (UnlockFunc, error) {

	activityCtx := workflow.WithLocalActivityOptions(ctx, workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	})

	var releaseLockChannelName string
	var execution workflow.Execution
	err := workflow.ExecuteLocalActivity(activityCtx, SignalWithStartMutexWorkflowActivity, s.lockNamespace,
		resourceID, s.currentWorkflowID, unlockTimeout).Get(ctx, &execution)
	if err != nil {
		return nil, err
	}

	isCanceled := false
	selector := workflow.NewSelector(ctx)
	selector.AddReceive(workflow.GetSignalChannel(ctx, AcquireLockSignalName), func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &releaseLockChannelName)
		workflow.GetLogger(ctx).Info("acquire lock", "lockName", releaseLockChannelName)
	})
	selector.AddReceive(ctx.Done(), func(c workflow.ReceiveChannel, more bool) {
		isCanceled = true
		workflow.GetLogger(ctx).Info("recieved cancel")
	})
	selector.Select(ctx)

	if !isCanceled {
		unlockFunc := func() error {
			return workflow.SignalExternalWorkflow(ctx, execution.ID, execution.RunID,
				releaseLockChannelName, s.currentWorkflowID).Get(ctx, nil)
		}
		return unlockFunc, nil
	}
	return nil, temporal.NewCanceledError()
}

// SignalWithStartMutexWorkflowActivity ...
func SignalWithStartMutexWorkflowActivity(
	ctx context.Context,
	namespace string,
	resourceID string,
	senderWorkflowID string,
	unlockTimeout time.Duration,
) (*workflow.Execution, error) {

	c := ctx.Value(ClientContextKey).(client.Client)
	workflowID := fmt.Sprintf(
		"%s:%s:%s",
		"mutex",
		namespace,
		resourceID,
	)
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: TaskQueue,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	wr, err := c.SignalWithStartWorkflow(
		ctx, workflowID, RequestLockSignalName, senderWorkflowID,
		workflowOptions, MutexWorkflowWithCancellation, namespace, resourceID, unlockTimeout, nil)

	if err != nil {
		activity.GetLogger(ctx).Error("Unable to signal with start workflow", "Error", err)
	} else {
		activity.GetLogger(ctx).Info("Signaled and started Workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())
	}

	return &workflow.Execution{
		ID:    wr.GetID(),
		RunID: wr.GetRunID(),
	}, nil
}

func SampleWorkflowWithMutex(
	ctx workflow.Context,
	resourceID string,
) error {

	currentWorkflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	logger := workflow.GetLogger(ctx)
	logger.Info("started", "currentWorkflowID", currentWorkflowID, "resourceID", resourceID)

	mutex := NewMutex(currentWorkflowID, "TestUseCase")
	unlockFunc, err := mutex.LockWithCancellation(ctx, resourceID, 1*time.Minute)
	if err != nil {
		return err
	}
	defer func() {
		if unlockFunc != nil {
			unlockFunc()
		}
	}()

	// emulate long running process
	logger.Info("critical operation started")
	_ = workflow.Sleep(ctx, 5*time.Second)
	logger.Info("critical operation finished")

	logger.Info("finished")
	return nil
}
