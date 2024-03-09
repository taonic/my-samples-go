package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/interceptor"
	sdkinterceptor "go.temporal.io/sdk/interceptor"
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
	var taskQueue = "my-task-queue" + uuid.New().String()
	w := worker.New(c, taskQueue, worker.Options{
		Interceptors: []sdkinterceptor.WorkerInterceptor{
			NewTimeoutInterceptor(3 * time.Second),
		},
	},
	)
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
	logger := workflow.GetLogger(ctx)
	defer func() {
		info := workflow.GetInfo(ctx)
		if errors.Is(ctx.Err(), workflow.ErrCanceled) {
			logger.Info("Add clean up code here")
		}
	}()

	for i := 0; i < 5; i++ {
		workflow.Sleep(ctx, 1*time.Second)
	}
	return nil
}

type timeoutInterceptor struct {
	interceptor.WorkerInterceptorBase
	duration time.Duration
}

func NewTimeoutInterceptor(duration time.Duration) interceptor.WorkerInterceptor {
	return &timeoutInterceptor{duration: duration}
}

func (w *timeoutInterceptor) InterceptActivity(
	ctx context.Context,
	next interceptor.ActivityInboundInterceptor,
) interceptor.ActivityInboundInterceptor {
	return next
}

func (w *timeoutInterceptor) InterceptWorkflow(
	ctx workflow.Context,
	next interceptor.WorkflowInboundInterceptor,
) interceptor.WorkflowInboundInterceptor {
	i := &workflowInboundInterceptor{root: w}
	i.Next = next
	return i
}

type workflowInboundInterceptor struct {
	interceptor.WorkflowInboundInterceptorBase
	root *timeoutInterceptor
}

func (t *workflowInboundInterceptor) ExecuteWorkflow(
	ctx workflow.Context,
	in *sdkinterceptor.ExecuteWorkflowInput,
) (interface{}, error) {
	ctxTmt := setTimeout(ctx, t.root.duration)
	ret, err := t.Next.ExecuteWorkflow(ctxTmt, in)
	return ret, err
}

func setTimeout(ctx workflow.Context, duration time.Duration) workflow.Context {
	timerCtx, timerCancel := workflow.WithCancel(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		workflow.Sleep(ctx, duration)
		timerCancel()
	})
	return timerCtx
}
