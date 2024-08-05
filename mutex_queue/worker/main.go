package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	mutex "github.com/taonic/my-samples-go/mutex_queue"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "mutex", worker.Options{
		BackgroundActivityContext: context.WithValue(context.Background(), mutex.ClientContextKey, c),
	})

	w.RegisterActivity(mutex.SignalWithStartMutexWorkflowActivity)
	w.RegisterWorkflow(mutex.MutexWorkflowWithCancellation)
	w.RegisterWorkflow(mutex.SampleWorkflowWithMutex)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
