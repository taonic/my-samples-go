Based on the below scenario:
https://community.temporal.io/t/mutex-workflow-and-how-to-track-workflow-that-have-requested-lock/12983

### Steps to run this sample:
1) Run a [Temporal service](https://github.com/temporalio/samples-go/tree/main/#how-to-use).
2) Run the following command to start the worker
```
go run mutex/worker/main.go
```
3) Run the following command to start the example
```
go run mutex/starter/main.go
```
