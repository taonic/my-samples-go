# QPS Sampling

Demonstrates MySQL QPS sampling using Temporal workflows and activities, with concurrent HelloWorld workflows.

## Prerequisites

Setup Temporal with Docker Compose using MySQL:

```bash
git clone https://github.com/temporalio/docker-compose.git
cd docker-compose
docker-compose -f docker-compose-mysql-es.yml up
```

## Running

1. Start the workers (in separate terminals):
```bash
# Terminal 1: Start QPS sampling worker
go run qps-worker/main.go

# Terminal 2: Start HelloWorld worker
go run hello-worker/main.go
```

2. Start the workflows:
```bash
# Start with default settings (1 HelloWorld workflow, 3 activities each)
go run starter/main.go

# Start with custom configuration
go run starter/main.go -workflows=5 -activities=10
```

## Configuration

### Command Line Flags
- `-workflows`: Number of concurrent HelloWorld workflows to start (default: 1)
- `-activities`: Number of activities per HelloWorld workflow (default: 3)

### Database Configuration
Update the MySQL connection string in `main.go` to match your database configuration.

### Enable MySQL Query Logs

To capture query logs in a table format, enable the general log:

```sql
SET GLOBAL general_log = 'ON';
SET GLOBAL log_output = 'TABLE';
```

Query logs will be stored in `mysql.general_log` table.

## Workflows

This sample runs two types of workflows on separate workers:

1. **QPS Sampling Workflow**: Continuously samples MySQL query logs to monitor database activity (runs on `qps-sampling` task queue)
2. **HelloWorld Workflow**: Executes a configurable number of simple greeting activities (runs on `hello-world` task queue)
