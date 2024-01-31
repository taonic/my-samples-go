package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	samples "github.com/taonic/my-samples-go"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	//channel := make(chan string)
	//go func() {
	//for i := 0; i < 50; i++ {
	//myi := i
	//var output string
	//output = fmt.Sprintf("%v", myi)
	//fmt.Printf("sending %v\n", output)
	//channel <- output
	//fmt.Printf("sent %v\n", output)
	//}
	//}()

	//go func() {
	//for i := 0; i < 50; i++ {
	//receivedValue := <-channel
	//fmt.Printf("received %v\n", receivedValue)
	//}
	//}()

	//time.Sleep(10 * time.Second)

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var clientOptions client.Options

	clientOptions, err := samples.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Println("Invalid arguments: %v. Connecting to local server.", err)
		clientOptions = client.Options{}
	}

	c, err := client.Dial(clientOptions)
	if err != nil {
		return err
	}
	defer c.Close()

	// Start worker
	log.Print("Running worker with limited concurrent activities")
	const taskQueue = "my-task-queue"
	w := worker.New(c, taskQueue, worker.Options{
		// Only allow so many at once
		MaxConcurrentActivityExecutionSize: 2,
	})
	w.RegisterWorkflow(MyWorkflow)
	w.RegisterActivity(MyActivity)
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	// Run workflow
	log.Print("Running workflow")
	run, err := c.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{TaskQueue: taskQueue},
		MyWorkflow,
	)
	if err != nil {
		return err
	} else if err := run.Get(context.Background(), nil); err != nil {
		return err
	}
	log.Printf("Workflow done")
	return nil
}

func MyActivity(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("hello %v", name), nil
}

func MyWorkflow(ctx workflow.Context) error {
	userChannel := workflow.NewChannel(ctx)

	for i := 0; i < 10; i++ {
		myi := i
		workflow.Go(ctx, func(gCtx workflow.Context) {
			gCtx = workflow.WithActivityOptions(gCtx,
				workflow.ActivityOptions{StartToCloseTimeout: 100 * time.Second},
			)
			var output string
			input := fmt.Sprintf("%v", myi)
			fmt.Printf("sending %v\n", input)
			result := workflow.ExecuteActivity(gCtx, MyActivity, input)
			result.Get(gCtx, &output)
			fmt.Printf("got %v\n", output)
			userChannel.Send(gCtx, output)
			fmt.Printf("sent %v\n", output)
		})
	}

	for i := 0; i < 10; i++ {
		var receivedValue string
		userChannel.Receive(ctx, &receivedValue)
		fmt.Printf("received %v\n", receivedValue)
	}
	_ = workflow.Await(ctx, func() bool {
		return true
	})

	return nil
}
