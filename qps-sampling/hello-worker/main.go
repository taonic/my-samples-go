package main

import (
	"log"
	"log/slog"
	"os"

	"go.temporal.io/sdk/client"
	tlog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"

	qps "github.com/taonic/my-samples-go/qps-sampling"
)

func main() {
	clientOptions := client.Options{
		Logger: tlog.NewStructuredLogger(
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			}))),
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, "hello-world", worker.Options{})
	w.RegisterWorkflow(qps.HelloWorldWorkflow)
	w.RegisterActivity(qps.HelloWorldActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
