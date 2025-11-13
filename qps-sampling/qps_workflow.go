package qpssampling

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

// QPSSample represents a QPS measurement
type QPSSample struct {
	Timestamp          time.Time
	QPS                float64
	Queries            []string
	LastEventTime      time.Time
	CompletedWorkflows int64
}

// QPSCountSample represents a QPS count measurement
type QPSCountSample struct {
	Timestamp          time.Time
	QueryCount         int64
	LastEventTime      time.Time
	CompletedWorkflows int64
}

// SampleDBQueriesActivity samples MySQL QPS
// func SampleDBQueriesActivity(ctx context.Context, lastSeen time.Time, lastCompletedCount int64) (QPSSample, error) {
// 	// Connect to MySQL (adjust connection string as needed)
// 	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/")
// 	if err != nil {
// 		return QPSSample{}, fmt.Errorf("failed to connect to MySQL: %w", err)
// 	}
// 	defer db.Close()

// 	// Get all arguments from general log since last event time
// 	var rows *sql.Rows
// 	if lastSeen.IsZero() {
// 		rows, err = db.Query("SELECT argument, event_time FROM mysql.general_log ORDER BY event_time DESC")
// 	} else {
// 		rows, err = db.Query("SELECT argument, event_time FROM mysql.general_log WHERE event_time > ? ORDER BY event_time DESC", lastSeen)
// 	}
// 	if err != nil {
// 		return QPSSample{}, fmt.Errorf("failed to query general log: %w", err)
// 	}
// 	defer rows.Close()

// 	var queries []string
// 	var maxEventTime time.Time
// 	for rows.Next() {
// 		var argument sql.NullString
// 		var eventTimeStr sql.NullString
// 		err := rows.Scan(&argument, &eventTimeStr)
// 		if err != nil {
// 			return QPSSample{}, fmt.Errorf("failed to scan row: %w", err)
// 		}
// 		if argument.Valid && eventTimeStr.Valid {
// 			arg := argument.String
// 			if len(arg) > 100 {
// 				arg = arg[:100]
// 			}
// 			queryWithTime := fmt.Sprintf("%s\t%s", eventTimeStr.String, arg)
// 			queries = append(queries, queryWithTime)
// 			if eventTime, err := time.Parse("2006-01-02 15:04:05", eventTimeStr.String); err == nil {
// 				if eventTime.After(maxEventTime) {
// 					maxEventTime = eventTime
// 				}
// 			}
// 		}
// 	}

// 	sample := QPSSample{
// 		Timestamp:     time.Now(),
// 		LastEventTime: maxEventTime,
// 	}

// 	// Count completed workflows
// 	c, err := client.Dial(client.Options{})
// 	if err == nil {
// 		resp, err := c.CountWorkflow(ctx, &workflowservice.CountWorkflowExecutionsRequest{
// 			Query: "ExecutionStatus='Completed'",
// 		})
// 		if err == nil {
// 			sample.CompletedWorkflows = resp.Count
// 			newCompleted := resp.Count - lastCompletedCount
// 			fmt.Printf("Completed workflows since last query: %d (total: %d)\n", newCompleted, resp.Count)
// 		}
// 		c.Close()
// 	}

// 	fmt.Printf("Queries count since last activity: %d\n", len(queries))
// 	fmt.Printf("=======================================================================\n")
// 	time.Sleep(1 * time.Second)
// 	return sample, nil
// }

// SampleDBQueryCountActivity counts DB queries and completed workflows
func SampleDBQueryCountActivity(ctx context.Context, lastSeen time.Time, lastCompletedCount int64, lastQueryCount int64) (QPSCountSample, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/")
	if err != nil {
		return QPSCountSample{}, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	defer db.Close()

	var variableName string
	var queryCount int64
	err = db.QueryRow("SHOW GLOBAL STATUS LIKE 'Queries'").Scan(&variableName, &queryCount)
	if err != nil {
		return QPSCountSample{}, fmt.Errorf("failed to get query count: %w", err)
	}

	sample := QPSCountSample{
		Timestamp:     time.Now(),
		QueryCount:    queryCount,
		LastEventTime: time.Now(),
	}

	c, err := client.Dial(client.Options{})
	if err == nil {
		resp, err := c.CountWorkflow(ctx, &workflowservice.CountWorkflowExecutionsRequest{
			Query: "ExecutionStatus='Completed'",
		})
		if err == nil {
			sample.CompletedWorkflows = resp.Count
			newCompleted := resp.Count - lastCompletedCount
			fmt.Printf("Completed workflows since last query: %d (total: %d)\n", newCompleted, resp.Count)
		}
		c.Close()
	}

	newQueries := queryCount - lastQueryCount
	fmt.Printf("DB Query count since last check: %d\n", newQueries)
	fmt.Printf("=======================================================================\n")
	time.Sleep(1 * time.Second)
	return sample, nil
}

// QPSSamplingWorkflow samples MySQL QPS periodically
func QPSSamplingWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var lastSeen time.Time
	var lastCompletedCount int64
	var lastQueryCount int64
	for {
		var sample QPSCountSample
		err := workflow.ExecuteActivity(ctx, SampleDBQueryCountActivity, lastSeen, lastCompletedCount, lastQueryCount).Get(ctx, &sample)
		if err == nil {
			lastSeen = sample.LastEventTime
			lastCompletedCount = sample.CompletedWorkflows
			lastQueryCount = sample.QueryCount
		}
		if err != nil {
			logger.Error("Failed to sample MySQL QPS", "error", err)
		}
	}
}
