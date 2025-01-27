# Temporal Go SDK Samples
A collection of samples demonstrating various features and patterns using the Temporal Go SDK.

## Samples

### [Bench Concurrent Workflow](/bench_concurrent_workflow)
Demonstrates workflow throughput optimization by partitioning concurrent executions into chunks.

### [Blocking Channel](/blocking_channel)
Shows patterns for handling blocking channels in workflows with error handling and timeouts.

### [Channels](/channels)
Illustrates communication between concurrent workflow executions using workflow channels.

### [Cloud API](/cloud-api)
Examples of managing Temporal Cloud resources using the Cloud API client.

### [Concurrent Cancel](/concurrent_cancel)
Shows handling of workflow cancellation with concurrent signals and activities.

### [CPU Intensive](/cpu_intensive)
Demonstrates management of CPU-intensive activities with concurrency controls.

### [Deadlock](/deadlock)
Illustrates deadlock detection and prevention in workflows.

### [FIFO Signals](/fifo_signals)
Shows implementation of First-In-First-Out signal processing in workflows.

### [GetHistory mTLS](/gethistorymtls)
Demonstrates retrieving workflow history using mTLS authentication.

### [List Workflows](/list_workflows)
Shows how to query and list workflows with various filters.

### [Max Concurrent Activities](/max_concurrent_activities)
Demonstrates limiting concurrent activity executions.

### [Mutex Queue](/mutex_queue)
Implements distributed mutex pattern using Temporal workflows.

### [OpenTelemetry](/opentelemetry)
Shows OpenTelemetry integration for tracing and metrics.

### [OTLP Metrics](/otlpmetrics)
Demonstrates metrics export using OpenTelemetry Protocol.

### [Parse Export Protobuf](/parse_export_protobuf)
Shows parsing of workflow history exported in protobuf format.

### [Query Schedules](/query_schedules)
Demonstrates querying workflow schedules using search attributes.

### [Race Context](/race_context)
Shows handling of race conditions with workflow context.

### [Race Query](/race-query)
Demonstrates safe handling of concurrent workflow queries.

### [Replay With Version And Marker](/replay_with_version_and_marker)
Shows workflow versioning and replay with markers.

### [Timeout Interceptor](/timeout_interceptor)
Implements custom timeout handling using workflow interceptors.

### [WaitGroup](/waitgroup)
Shows coordination of multiple concurrent workflows using WaitGroup pattern.

## Prerequisites
- Go 1.21 or later
- Temporal Server running locally or a Temporal Cloud account

## Building and Running
Each sample can be run using:

```bash
go run <sample_directory>/main.go
```
