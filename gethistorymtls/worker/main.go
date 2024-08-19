package main

import (
	"context"
	"log"
	"os"

	"github.com/taonic/my-samples-go/gethistorymtls"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	clientOptions, err := gethistorymtls.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	ctx := context.WithValue(context.Background(), gethistorymtls.TemporalClientKey, c)

	w := worker.New(c, "hello-world-mtls", worker.Options{
		BackgroundActivityContext: ctx,
	})

	w.RegisterWorkflow(gethistorymtls.Workflow)
	w.RegisterActivity(gethistorymtls.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
