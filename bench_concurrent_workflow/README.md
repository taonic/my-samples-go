### Benchmark concurrent Workflow execution

This example demonstrates the effectiveness of partitioning concurrent Workflow executions into chunks
as a standard strategy to workaround the individual Workflow's throughput limitation. Due to the consistency requirement,
Temporal server needs to put a lock on the Workflow while performing transaction, which may cause lock contention when
large amount of transactions are queued up.

Breaking down 200 Workflows into various partitions shows the below result:
```
benchmark/partitions:1/partitionSize:200/activities:1 was done in 7.895454334s
benchmark/partitions:2/partitionSize:100/activities:1 was done in 3.042504292s
benchmark/partitions:4/partitionSize:50/activities:1 was done in 1.722763292s
benchmark/partitions:8/partitionSize:25/activities:1 was done in 1.207211083s
benchmark/partitions:20/partitionSize:10/activities:1 was done in 1.036813291s
benchmark/partitions:40/partitionSize:5/activities:1 was done in 1.422771375s
```

To run the sample:
```
go run bench_concurrent_workflow/main.go
```
