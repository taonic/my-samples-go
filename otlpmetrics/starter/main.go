package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	metrics "github.com/taonic/my-samples-go/otlpmetrics"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client.", err)
	}
	defer c.Close()

	for true {
		workflowOptions := client.StartWorkflowOptions{
			ID:        "metrics_workflowID",
			TaskQueue: "metrics",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, metrics.Workflow)
		if err != nil {
			log.Fatalln("Unable to execute workflow.", err)
		}

		log.Println("Started workflow.", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

		// Synchronously wait for the workflow completion.
		err = we.Get(context.Background(), nil)
		if err != nil {
			log.Fatalln("Unable to wait for workflow completition.", err)
		}
		time.Sleep(1 * time.Second)
	}

	log.Println("Check metrics at http://localhost:9090/metrics")
}
