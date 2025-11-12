package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Start workflow with 5-second sampling interval
	workflowOptions := client.StartWorkflowOptions{
		ID:        "qps-sampling-workflow",
		TaskQueue: "qps-sampling",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, "QPSSamplingWorkflow", 5)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}