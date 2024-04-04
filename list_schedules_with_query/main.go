package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/taonic/my-samples-go/lib"
	workflowpb "go.temporal.io/api/workflow/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
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

	projectID := "12345"

	// Create schedules
	for i := 0; i < 2; i++ {
		createSchedule(ctx, c, projectID)
	}

	time.Sleep(2 * time.Second) // Wait for Workflows to be indexed.

	// Query schedules
	// Note, this query is using attributes associate with internal implementation details that are subject to change.
	query := fmt.Sprintf("ExecutionStatus='Running' and ProjectID='%s' and TemporalNamespaceDivision = 'TemporalScheduler'", projectID)
	workflowExecutions, err := GetWorkflows(ctx, c, query)
	if err != nil {
		log.Fatalln("Error listing schedules", err)
	}

	// Describe schedules
	for i, wf := range workflowExecutions {
		wfID := wf.Execution.WorkflowId
		scheduleID := strings.Split(wfID, "temporal-sys-scheduler:")[1]
		desc, err := c.ScheduleClient().GetHandle(ctx, scheduleID).Describe(ctx)
		if err != nil {
			log.Fatalln("Failed to describe schedule", err)
		}
		fmt.Printf("Schedule %d: %+v \n", i, desc)
	}

	return nil
}

func createSchedule(ctx context.Context, c client.Client, projectID string) error {
	scheduleID := "schedule_" + uuid.New().String()
	workflowID := "schedule_workflow_" + uuid.New().String()
	_, err := c.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID:   scheduleID,
		Spec: client.ScheduleSpec{},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  "SampleScheduleWorkflow",
			TaskQueue: "schedule",
		},
		SearchAttributes: map[string]interface{}{
			"ProjectID": projectID,
		},
	})
	if err != nil {
		return err
	}

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
