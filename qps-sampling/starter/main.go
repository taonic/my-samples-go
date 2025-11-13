package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	qps "github.com/taonic/my-samples-go/qps-sampling"
)

func main() {
	workflowsPerSecond := flag.Int("workflows", 10, "Number of workflows to start per second")
	activitiesPerWorkflow := flag.Int("activities", 5, "Number of activities per workflow")
	flag.Parse()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Start QPS sampling workflow
	qpsOptions := client.StartWorkflowOptions{
		ID:        "qps-sampling-workflow",
		TaskQueue: "qps-sampling",
	}
	_, err = c.ExecuteWorkflow(context.Background(), qpsOptions, qps.QPSSamplingWorkflow)
	if err != nil {
		log.Printf("Failed to start QPS sampling workflow: %v", err)
	} else {
		log.Println("Started QPS sampling workflow")
	}

	startConcurrentWorkflows(c, *workflowsPerSecond, *activitiesPerWorkflow)
}

func startConcurrentWorkflows(c client.Client, workflowsPerSecond, activitiesPerWorkflow int) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	workflowCount := 0
	for range ticker.C {
		for i := 0; i < workflowsPerSecond; i++ {
			go func(id int) {
				workflowID := fmt.Sprintf("hello-workflow-%d-%d", time.Now().Unix(), id)
				options := client.StartWorkflowOptions{
					ID:        workflowID,
					TaskQueue: "hello-world",
				}
				_, err := c.ExecuteWorkflow(context.Background(), options, qps.HelloWorldWorkflow, activitiesPerWorkflow)
				if err != nil {
					log.Printf("Failed to start workflow %s: %v", workflowID, err)
				}
			}(i)
		}
		workflowCount += workflowsPerSecond
		log.Printf("Started %d concurrent workflows (total: %d)", workflowsPerSecond, workflowCount)
	}
}
