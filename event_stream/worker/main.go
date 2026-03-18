package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	es "github.com/taonic/my-samples-go/event_stream"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, es.TaskQueue, worker.Options{})
	w.RegisterWorkflow(es.MoneyTransferWorkflow)
	w.RegisterActivity(es.ValidateTransfer)
	w.RegisterActivity(es.Withdraw)
	w.RegisterActivity(es.Deposit)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
