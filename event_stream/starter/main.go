package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"

	es "github.com/taonic/my-samples-go/event_stream"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowID := "money-transfer-" + uuid.NewString()[:8]
	run, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: es.TaskQueue,
	}, es.MoneyTransferWorkflow, es.TransferInput{
		FromAccount: "ACC-001",
		ToAccount:   "ACC-002",
		Amount:      250.00,
		ReferenceID: workflowID,
	})
	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}
	log.Printf("Started workflow %s (run %s)\n", run.GetID(), run.GetRunID())

	// Long-poll workflow history for SideEffect marker events
	ctx := context.Background()
	iter := c.GetWorkflowHistory(ctx, run.GetID(), run.GetRunID(), true, enumspb.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)

	dc := converter.GetDefaultDataConverter()
	for iter.HasNext() {
		event, err := iter.Next()
		if err != nil {
			log.Fatalln("Error reading history", err)
		}

		attrs := event.GetMarkerRecordedEventAttributes()
		if attrs == nil || attrs.GetMarkerName() != "SideEffect" {
			continue
		}

		dataPayloads, ok := attrs.GetDetails()["data"]
		if !ok || len(dataPayloads.GetPayloads()) == 0 {
			continue
		}

		var te es.TransferEvent
		if err := dc.FromPayload(dataPayloads.GetPayloads()[0], &te); err != nil {
			log.Printf("Failed to decode event: %v", err)
			continue
		}

		pretty, _ := json.MarshalIndent(te, "", "  ")
		fmt.Printf("📌 Event: %s\n", pretty)
	}

	log.Println("Workflow completed")
}
