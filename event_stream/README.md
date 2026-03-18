# Event Stream

Demonstrates using `workflow.SideEffect` to emit custom events during workflow execution, and using `GetWorkflowHistory` (long-poll) on the client side to stream and decode those events in real time.

Uses a money transfer workflow as an example — after each activity (validate, withdraw, deposit), a `TransferEvent` is recorded via `SideEffect`. The starter polls the workflow history for `MarkerRecorded` events with marker name `"SideEffect"`, decodes the payload, and prints it.

## Sequence Diagram

```mermaid
sequenceDiagram
    participant C as Client
    participant T as Temporal Server
    participant W as Workflow

    C->>T: StartWorkflow
    C->>T: GetWorkflowHistory (long-poll)

    W->>T: SideEffect(transfer_initiated)
    T-->>C: TransferEvent(transfer_initiated)
    W->>W: Activity: ValidateTransfer
    W->>T: SideEffect(validation completed)
    T-->>C: TransferEvent(validation completed)
    W->>W: Activity: Withdraw
    W->>T: SideEffect(withdraw completed)
    T-->>C: TransferEvent(withdraw completed)
    W->>W: Activity: Deposit
    W->>T: SideEffect(deposit completed)
    T-->>C: TransferEvent(deposit completed)
    W->>T: SideEffect(transfer_completed)
    T-->>C: TransferEvent(transfer_completed)
```

## Running

Start the worker:
```bash
go run event_stream/worker/main.go
```

In another terminal, start the workflow and stream events:
```bash
go run event_stream/starter/main.go
```
