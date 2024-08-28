package opentelemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

var tracer trace.Tracer

func init() {
	// Name the tracer after the package, or the service if you are in main
	tracer = otel.Tracer("github.com/temporalio/samples-go/otel")
}

func Workflow(ctx workflow.Context, name string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})

	err := workflow.ExecuteActivity(ctx, Activity).Get(ctx, nil)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return err
	}

	span := ctx.Value("my-span-context-key").(trace.Span)
	sCtx := trace.ContextWithSpanContext(context.Background(), span.SpanContext())
	_, signalSpan := tracer.Start(sCtx, "my-cross-activity-span", trace.WithTimestamp(workflow.Now(ctx)))

	err = workflow.ExecuteActivity(ctx, Activity).Get(ctx, nil)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return err
	}

	fmt.Println("You have 5s to restart the worker to force a replay")
	workflow.Sleep(ctx, 5*time.Second)

	err = workflow.ExecuteActivity(ctx, Activity).Get(ctx, nil)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return err
	}

	signalSpan.End(trace.WithTimestamp(workflow.Now(ctx)))

	logger.Info("HelloWorld workflow completed.")
	return nil
}

func Activity(ctx context.Context, name string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)

	// Get current span and add new attributes
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Bool("isTrue", true), attribute.String("stringAttr", "Ciao"))

	// Create a child span
	_, childSpan := tracer.Start(ctx, "custom-span")
	time.Sleep(1 * time.Second)
	childSpan.End()

	time.Sleep(1 * time.Second)

	// Add an event to the current span
	span.AddEvent("Done Activity")

	return nil
}
