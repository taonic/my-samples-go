package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/cloud-sdk/cloudclient"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("TEMPORAL_CLOUD_API_KEY")
	if apiKey == "" {
		log.Fatal("TEMPORAL_CLOUD_API_KEY environment variable is required")
	}

	// Create cloud client
	cloudClient, err := cloudclient.New(cloudclient.Options{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create cloud client: %v", err)
	}
	defer cloudClient.Close()

	ctx := context.Background()

	// Get usage for the last 30 days
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -30)

	req := &cloudservice.GetUsageRequest{
		StartTimeInclusive: timestamppb.New(startTime),
		EndTimeExclusive:   timestamppb.New(endTime),
	}

	resp, err := cloudClient.CloudService().GetUsage(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get usage: %v", err)
	}

	//fmt.Printf("Usage data from %s to %s:\n", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
	fmt.Printf("Summary: %+v\n", resp)

	//if len(resp.UsageByNamespace) > 0 {
	//fmt.Println("\nUsage by namespace:")
	//for _, usage := range resp.UsageByNamespace {
	//fmt.Printf("  Namespace: %s\n", usage.Namespace)
	//fmt.Printf("    Actions: %d\n", usage.Summary.ActionCount)
	//fmt.Printf("    Storage: %d bytes\n", usage.Summary.StorageBytes)
	//}
	//}
}
