## Concurrent Cancellation

This sample is adapted from https://github.com/temporalio/samples-go/tree/main/cancellation to demonstrate how cancellation can be handle concurrently with blocking Signals and Activities

<!-- @@@SNIPSTART samples-go-cancellation-readme -->
Make sure the [Temporal Server is running locally](https://docs.temporal.io/application-development/foundations#run-a-development-cluster).

From the root of the project, start a Worker:

```bash
go run cancellation/worker/main.go
```

Start the Workflow Execution:

```bash
go run cancellation/starter/main.go
```

Cancel the Workflow Execution:

```bash
go run cancellation/cancel/main.go
```
<!-- @@@SNIPEND -->
