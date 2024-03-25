### Steps to run this sample:

This example show cases an alternative way to set Workflow execution timeout that can be gracefully handled, as opposed to the server side timeout (WorkflowOptions.WorkflowExecutionTimeout).

It uses the inbound interceptor `ExecuteWorkflow` to set a timer at the beginning of each Workflow execution and cancel the Workflow context when the timer fires.


```
go run timeout_interceptor/main.go
```
