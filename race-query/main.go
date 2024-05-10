package main

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var (
	c          client.Client
	wg         *sync.WaitGroup
	workflowID = uuid.New().String()
	taskQueue  = "my-task-queue" + uuid.New().String()
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error

	c, err = client.Dial(client.Options{})
	if err != nil {
		return err
	}
	defer c.Close()

	// Start worker
	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(QueryWorkflow)
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	wg = &sync.WaitGroup{}

	// Start WF
	log.Print("Running workflow")
	_, err = c.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        workflowID,
			TaskQueue: taskQueue,
		},
		QueryWorkflow,
		0,
	)
	if err != nil {
		return err
	}

	// Repeat signal and query concurrently
	signal(ctx)
	query(ctx)

	wg.Wait()

	return nil
}

func signal(ctx context.Context) {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			err := c.SignalWorkflow(ctx, workflowID, "", "Signal", nil)
			if err != nil {
				log.Println("Failed to signal", err)
			}
		}()
	}
}

func query(ctx context.Context) {
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := c.QueryWorkflow(context.Background(), workflowID, "", "state")
			if err != nil {
				log.Println("Unable to query workflow", "Error", err)
				if err.Error() == "Workflow task is not scheduled yet." {
					log.Fatalln("task not scheduled error found")
				}
				return
			}
			var result interface{}
			if err := resp.Get(&result); err != nil {
				log.Fatalln("Unable to decode query result", err)
			}
			log.Println("Query succeed with RunID", result)
		}()
	}
}

func QueryWorkflow(ctx workflow.Context, runCount int) error {
	sigCounter := 0

	logger := workflow.GetLogger(ctx)
	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (string, error) {
		return workflow.GetInfo(ctx).WorkflowExecution.RunID, nil
	})
	if err != nil {
		return err
	}

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(workflow.GetSignalChannel(ctx, "Signal"), func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		sigCounter++
	})

	for sigCounter < 5 || selector.HasPending() {
		selector.Select(ctx)
	}

	if runCount == 10 {
		return nil
	}

	logger.Info("About to continue as new")

	return workflow.NewContinueAsNewError(ctx, QueryWorkflow, runCount+1)
}
