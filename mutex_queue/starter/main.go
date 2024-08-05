package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"github.com/google/uuid"
	mutex "github.com/taonic/my-samples-go/mutex_queue"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	taskQueue := "mutex"

	// This workflow ID can be user business logic identifier as well.
	resourceID := "mutex_resource"
	workflow1Options := client.StartWorkflowOptions{
		ID:        "A_" + uuid.New().String(),
		TaskQueue: taskQueue,
	}

	workflow2Options := client.StartWorkflowOptions{
		ID:        "B_" + uuid.New().String(),
		TaskQueue: taskQueue,
	}

	workflow3Options := client.StartWorkflowOptions{
		ID:        "C_" + uuid.New().String(),
		TaskQueue: taskQueue,
	}

	workflow4Options := client.StartWorkflowOptions{
		ID:        "D_" + uuid.New().String(),
		TaskQueue: taskQueue,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflow1Options, mutex.SampleWorkflowWithMutex, resourceID)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	} else {
		log.Println(fmt.Sprintf("Started workflow [%s]", we.GetID()))
	}
	time.Sleep(100 * time.Millisecond)

	we, err = c.ExecuteWorkflow(context.Background(), workflow2Options, mutex.SampleWorkflowWithMutex, resourceID)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	} else {
		log.Println(fmt.Sprintf("Started workflow [%s]", we.GetID()))
	}
	time.Sleep(100 * time.Millisecond)

	we, err = c.ExecuteWorkflow(context.Background(), workflow3Options, mutex.SampleWorkflowWithMutex, resourceID)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	} else {
		log.Println(fmt.Sprintf("Started workflow [%s]", we.GetID()))
	}
	time.Sleep(100 * time.Millisecond)

	we, err = c.ExecuteWorkflow(context.Background(), workflow4Options, mutex.SampleWorkflowWithMutex, resourceID)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	} else {
		log.Println(fmt.Sprintf("Started workflow [%s]", we.GetID()))
	}
}
