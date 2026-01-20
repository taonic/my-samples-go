package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"github.com/taonic/my-samples-go/lib"
)

const taskQueue = "schedule-update-rps-queue"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	clientOptions, err := lib.ParseClientOptionFlags(os.Args[1:])
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
	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(DummyWorkflow)
	if err := w.Start(); err != nil {
		return err
	}
	defer w.Stop()

	numSchedules := 200
	scheduleIDs := make([]string, numSchedules)

	// Create schedules
	log.Printf("Creating %d schedules...", numSchedules)
	for i := 0; i < numSchedules; i++ {
		scheduleID := fmt.Sprintf("schedule-update-test-%d-%d", time.Now().Unix(), i)
		scheduleIDs[i] = scheduleID

		_, err := c.ScheduleClient().Create(ctx, client.ScheduleOptions{
			ID: scheduleID,
			Spec: client.ScheduleSpec{
				Intervals: []client.ScheduleIntervalSpec{
					{Every: 1 * time.Hour},
				},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        fmt.Sprintf("workflow-%s", scheduleID),
				Workflow:  DummyWorkflow,
				TaskQueue: taskQueue,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create schedule %s: %w", scheduleID, err)
		}
	}
	log.Printf("Created %d schedules", numSchedules)

	// Update schedules concurrently and measure RPS
	log.Println("Starting schedule updates...")
	startTime := time.Now()
	var wg sync.WaitGroup
	updateCount := 0
	var mu sync.Mutex

	for _, scheduleID := range scheduleIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			handle := c.ScheduleClient().GetHandle(ctx, id)
			err := handle.Update(ctx, client.ScheduleUpdateOptions{
				DoUpdate: func(input client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
					input.Description.Schedule.Spec.Intervals = []client.ScheduleIntervalSpec{
						{Every: 30 * time.Minute},
					}
					return &client.ScheduleUpdate{Schedule: &input.Description.Schedule}, nil
				},
			})
			if err != nil {
				log.Printf("Failed to update schedule %s: %v", id, err)
				return
			}
			mu.Lock()
			updateCount++
			mu.Unlock()
		}(scheduleID)
	}

	wg.Wait()
	duration := time.Since(startTime)
	rps := float64(updateCount) / duration.Seconds()

	log.Printf("\n=== Results ===")
	log.Printf("Total schedules updated: %d", updateCount)
	log.Printf("Duration: %v", duration)
	log.Printf("RPS: %.2f updates/second", rps)

	// Cleanup
	log.Println("\nCleaning up schedules...")
	for _, scheduleID := range scheduleIDs {
		handle := c.ScheduleClient().GetHandle(ctx, scheduleID)
		_ = handle.Delete(ctx)
	}

	return nil
}

func DummyWorkflow(ctx workflow.Context) error {
	return nil
}
