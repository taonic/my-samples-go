Based on the below scenario:
https://community.temporal.io/t/sharing-discovery-setup-temporal-with-nx-monorepo-and-nest-js/6696/5

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
