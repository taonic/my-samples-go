package main

import (
	"fmt"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	// TODO: change path to your generated proto
	"github.com/taonic/my-samples-go/parse_export_protobuf/export"

	"go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"

	// TODO: change path to temporal repo
	ossserialization "go.temporal.io/server/common/persistence/serialization"
)

func extractWorkflowHistoriesFromFile(filename string) ([]*export.Workflow, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading from file: %v", err)
	}

	blob := &common.DataBlob{
		EncodingType: enumspb.ENCODING_TYPE_PROTO3,
		Data:         bytes,
	}

	result := &export.ExportedWorkflows{}

	err = ossserialization.ProtoDecodeBlob(blob, result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file: %w", err)
	}

	workflows := result.Workflows
	for _, workflow := range workflows {
		history := workflow.History
		if history == nil {
			return nil, fmt.Errorf("history is nil")
		}
	}

	return workflows, nil
}

func printWorkflow(workflow *export.Workflow) {
	// Pretty print the workflow
	marshaler := jsonpb.Marshaler{
		Indent:       "\t",
		EmitDefaults: true,
	}

	str, err := marshaler.MarshalToString(workflow.History)
	if err != nil {
		fmt.Println("error print workflow history: ", err)
		os.Exit(1)
	}
	print(str)
}

func printAllWorkflows(workflowHistories []*export.Workflow) {
	for _, workflow := range workflowHistories {
		printWorkflow(workflow)
	}
}

func printWorkflowHistory(workflowID string, workflowHistories []*export.Workflow) {
	if workflowID == "" {
		fmt.Println("invalid workflow ID")
		os.Exit(1)
	}

	for _, workflow := range workflowHistories {
		if workflow.History.Events[0].GetWorkflowExecutionStartedEventAttributes().WorkflowId == workflowID {
			fmt.Println("Printing workflow history for workflow ID: ", workflowID)
			printWorkflow(workflow)
		}
	}

	fmt.Println("No workflow found with workflow ID: ", workflowID)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a path to a file")
		os.Exit(1)
	}

	filename := os.Args[1]

	fmt.Println("Deserializing export workflow history from file: ", filename)

	workflowHistories, err := extractWorkflowHistoriesFromFile(filename)

	if err != nil {
		fmt.Println("error extracting workflow histories: ", err)
		os.Exit(1)
	}

	fmt.Println("Successfully deserialized workflow histories")
	fmt.Println("Total number of workflow histories: ", len(workflowHistories))

	fmt.Println("Choose an option:")
	fmt.Println("1. Print out all the workflows")
	fmt.Println("2. Print out the workflow hisotry of a specific workflow. Enter the workflow")
}
