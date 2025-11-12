package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.temporal.io/sdk/client"
	tlog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// QPSSample represents a QPS measurement
type QPSSample struct {
	Timestamp     time.Time
	QPS           float64
	Queries       []string
	LastEventTime time.Time
}

// SampleMySQLQPSActivity samples MySQL QPS
func SampleMySQLQPSActivity(ctx context.Context, lastSeen time.Time) (QPSSample, error) {
	// Connect to MySQL (adjust connection string as needed)
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/")
	if err != nil {
		return QPSSample{}, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	defer db.Close()

	// Get all arguments from general log since last event time
	var rows *sql.Rows
	if lastSeen.IsZero() {
		rows, err = db.Query("SELECT argument, event_time FROM mysql.general_log ORDER BY event_time DESC")
	} else {
		rows, err = db.Query("SELECT argument, event_time FROM mysql.general_log WHERE event_time > ? ORDER BY event_time DESC", lastSeen)
	}
	if err != nil {
		return QPSSample{}, fmt.Errorf("failed to query general log: %w", err)
	}
	defer rows.Close()

	var queries []string
	var maxEventTime time.Time
	for rows.Next() {
		var argument sql.NullString
		var eventTimeStr sql.NullString
		err := rows.Scan(&argument, &eventTimeStr)
		if err != nil {
			return QPSSample{}, fmt.Errorf("failed to scan row: %w", err)
		}
		if argument.Valid && eventTimeStr.Valid {
			arg := argument.String
			if len(arg) > 100 {
				arg = arg[:100]
			}
			queryWithTime := fmt.Sprintf("%s\t%s", eventTimeStr.String, arg)
			queries = append(queries, queryWithTime)
			if eventTime, err := time.Parse("2006-01-02 15:04:05", eventTimeStr.String); err == nil {
				if eventTime.After(maxEventTime) {
					maxEventTime = eventTime
				}
			}
		}
	}

	sample := QPSSample{
		Timestamp:     time.Now(),
		LastEventTime: maxEventTime,
	}

	fmt.Printf("Queries:\n%s\n", strings.Join(queries, "\n"))
	fmt.Printf("=======================================================================\n")
	fmt.Printf("Queries count since last activity: %d\n", len(queries))
	fmt.Printf("=======================================================================\n")
	time.Sleep(2 * time.Second)
	return sample, nil
}

// QPSSamplingWorkflow samples MySQL QPS periodically
func QPSSamplingWorkflow(ctx workflow.Context, intervalSeconds int) error {
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var lastSeen time.Time
	for {
		var sample QPSSample
		err := workflow.ExecuteActivity(ctx, SampleMySQLQPSActivity, lastSeen).Get(ctx, &sample)
		if err == nil {
			lastSeen = sample.LastEventTime
		}
		if err != nil {
			logger.Error("Failed to sample MySQL QPS", "error", err)
		}
	}
}

func main() {
	// Create the client object just once per process
	clientOptions := client.Options{
		Logger: tlog.NewStructuredLogger(
			slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo, // Set your desired log level here
			}))),
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, "qps-sampling", worker.Options{})
	w.RegisterWorkflow(QPSSamplingWorkflow)
	w.RegisterActivity(SampleMySQLQPSActivity)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
