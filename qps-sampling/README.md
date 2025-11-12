# QPS Sampling

Demonstrates MySQL QPS sampling using Temporal workflows and activities.

## Prerequisites

Setup Temporal with Docker Compose using MySQL:

```bash
git clone https://github.com/temporalio/docker-compose.git
cd docker-compose
docker-compose -f docker-compose-mysql-es.yml up
```

## Running

1. Start the worker:
```bash
go run main.go
```

2. Start the workflow:
```bash
go run starter/main.go
```

## Configuration

Update the MySQL connection string in `main.go` to match your database configuration.

### Enable MySQL Query Logs

To capture query logs in a table format, enable the general log:

```sql
SET GLOBAL general_log = 'ON';
SET GLOBAL log_output = 'TABLE';
```

Query logs will be stored in `mysql.general_log` table.
