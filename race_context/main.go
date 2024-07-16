package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

/*
This sample demonstrates sharing workflow context leads to:
Error getState: illegal access from outside of workflow context StackTrace coroutine root [panic]
*/

var (
	c         client.Client
	taskQueue = "my-task-queue" + uuid.New().String()
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
	wf := MyWorkflow{}
	w.RegisterWorkflow(wf.RaceWorkflow) // Don't do this. The method receiver will leak state between workflows.
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	// Run workflow
	log.Print("Running workflow concurrently")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			run, err := c.ExecuteWorkflow(
				ctx,
				client.StartWorkflowOptions{TaskQueue: taskQueue},
				wf.RaceWorkflow,
			)
			if err != nil {
				panic(err)
			} else if err := run.Get(context.Background(), nil); err != nil {
				panic(err)
			}
			log.Printf("Workflow done")
		}()
	}
	wg.Wait()
	return nil
}

type MyWorkflow struct {
	ctx workflow.Context
}

func (w *MyWorkflow) RaceWorkflow(ctx workflow.Context) error {
	w.ctx = ctx // Please don't do this as ctx will be leaked to other workflows.
	workflow.Sleep(w.ctx, 1*time.Second)
	workflow.Sleep(w.ctx, 1*time.Second)
	return nil
}
