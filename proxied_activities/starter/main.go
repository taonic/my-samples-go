package main

import (
	"context"
	"fmt"
	"log"

	proxied "github.com/taonic/my-samples-go/proxied_activities"
	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	for i := 0; i < 10; i++ {
		workflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("hello-proxy-workflow-%d", i),
			TaskQueue: "hello-proxy",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, proxied.HelloProxyWorkflow, "World")
		if err != nil {
			log.Printf("Unable to execute workflow %d: %v", i, err)
			continue
		}

		var result string
		err = we.Get(context.Background(), &result)
		if err != nil {
			log.Printf("Unable to get workflow %d result: %v", i, err)
			continue
		}

		log.Printf("Workflow %d result: %s", i, result)
	}
}
