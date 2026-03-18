package event_stream

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "event-stream-queue"

type TransferEvent struct {
	Step      string    `json:"step"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount,omitempty"`
	Account   string    `json:"account,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message,omitempty"`
}

type TransferInput struct {
	FromAccount string
	ToAccount   string
	Amount      float64
	ReferenceID string
}

func MoneyTransferWorkflow(ctx workflow.Context, input TransferInput) error {
	ao := workflow.ActivityOptions{StartToCloseTimeout: 30 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, ao)

	emitEvent := func(event TransferEvent) {
		workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
			return event
		})
	}

	emitEvent(TransferEvent{
		Step: "transfer_initiated", Status: "started",
		Amount: input.Amount, Message: fmt.Sprintf("Transfer %s -> %s", input.FromAccount, input.ToAccount),
		Timestamp: workflow.Now(ctx),
	})

	// Validate
	var valid bool
	if err := workflow.ExecuteActivity(ctx, ValidateTransfer, input).Get(ctx, &valid); err != nil {
		return err
	}
	emitEvent(TransferEvent{
		Step: "validation", Status: "completed", Timestamp: workflow.Now(ctx),
		Message: fmt.Sprintf("Validated transfer of $%.2f", input.Amount),
	})

	// Withdraw
	if err := workflow.ExecuteActivity(ctx, Withdraw, input.FromAccount, input.Amount).Get(ctx, nil); err != nil {
		return err
	}
	emitEvent(TransferEvent{
		Step: "withdraw", Status: "completed", Account: input.FromAccount,
		Amount: input.Amount, Timestamp: workflow.Now(ctx),
	})

	// Deposit
	if err := workflow.ExecuteActivity(ctx, Deposit, input.ToAccount, input.Amount).Get(ctx, nil); err != nil {
		return err
	}
	emitEvent(TransferEvent{
		Step: "deposit", Status: "completed", Account: input.ToAccount,
		Amount: input.Amount, Timestamp: workflow.Now(ctx),
	})

	emitEvent(TransferEvent{
		Step: "transfer_completed", Status: "success", Timestamp: workflow.Now(ctx),
		Message: fmt.Sprintf("Successfully transferred $%.2f from %s to %s", input.Amount, input.FromAccount, input.ToAccount),
	})

	return nil
}

func ValidateTransfer(ctx context.Context, input TransferInput) (bool, error) {
	activity.GetLogger(ctx).Info("Validating transfer", "from", input.FromAccount, "to", input.ToAccount, "amount", input.Amount)
	time.Sleep(500 * time.Millisecond)
	return true, nil
}

func Withdraw(ctx context.Context, account string, amount float64) error {
	activity.GetLogger(ctx).Info("Withdrawing", "account", account, "amount", amount)
	time.Sleep(1 * time.Second)
	return nil
}

func Deposit(ctx context.Context, account string, amount float64) error {
	activity.GetLogger(ctx).Info("Depositing", "account", account, "amount", amount)
	time.Sleep(1 * time.Second)
	return nil
}
