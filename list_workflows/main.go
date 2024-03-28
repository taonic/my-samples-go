package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/taonic/my-samples-go/lib"
	workflowpb "go.temporal.io/api/workflow/v1"
	"go.temporal.io/api/workflowservice/v1"
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
	clientOptions, err := lib.ParseClientOptionFlags(os.Args[1:])

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
	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(MyWorkflow)
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	// Run workflow
	log.Print("Running workflow")
	prefix := fmt.Sprintf("MyWorkflow-%s", uuid.New().String())
	for i := 0; i < 10; i++ {
		run, err := c.ExecuteWorkflow(
			ctx,
			client.StartWorkflowOptions{ID: fmt.Sprintf("%s-%d", prefix, i), TaskQueue: taskQueue},
			MyWorkflow,
		)
		if err != nil {
			return err
		} else if err := run.Get(context.Background(), nil); err != nil {
			return err
		}
	}

	time.Sleep(2 * time.Second) // Wait for Workflows to be indexed.

	// Query workflows
	query := fmt.Sprintf("WorkflowId between '%s-5' and '%s-8'", prefix, prefix)
	workflowExecutions, err := GetWorkflows(ctx, c, query)
	if err != nil {
		log.Fatalln("Error listing workflows", err)
	}
	log.Println("Should find 3 workflows", "Found", len(workflowExecutions))

	return nil
}

// GetWorkflows calls ListWorkflow with query and gets all workflow exection infos in a list.
func GetWorkflows(ctx context.Context, c client.Client, query string) ([]*workflowpb.WorkflowExecutionInfo, error) {
	var nextPageToken []byte
	var workflowExecutions []*workflowpb.WorkflowExecutionInfo
	for {
		resp, err := c.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
			Query:         query,
			NextPageToken: nextPageToken,
		})
		if err != nil {
			return nil, err
		}
		workflowExecutions = append(workflowExecutions, resp.Executions...)
		nextPageToken = resp.NextPageToken
		if len(nextPageToken) == 0 {
			return workflowExecutions, nil
		}
	}
}

func MyWorkflow(ctx workflow.Context) (string, error) {
	return "", nil
}
