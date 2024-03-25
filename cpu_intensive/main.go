package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/google/uuid"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/sync/errgroup"
)

func main() {
	runtime.GOMAXPROCS(4)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := client.Dial(client.Options{
		MetricsHandler: sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
			ListenAddress: "0.0.0.0:9091",
			TimerType:     "histogram",
		})),
	})
	if err != nil {
		return err
	}
	defer c.Close()

	var taskQueue = "my-task-queue" + uuid.New().String()

	// Start worker
	var group errgroup.Group
	work := func() error {
		w := worker.New(c, taskQueue, worker.Options{
			MaxConcurrentActivityTaskPollers: 100,
		})
		w.RegisterWorkflow(MyWorkflow)
		w.RegisterActivity(MyActivity)
		err = w.Run(worker.InterruptCh())
		if err != nil {
			return err
		}
		return nil
	}
	group.Go(work)

	// Run workflow
	log.Print("Running workflow")
	run, err := c.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{TaskQueue: taskQueue},
		MyWorkflow,
	)
	if err != nil {
		return err
	} else if err := run.Get(context.Background(), nil); err != nil {
		return err
	}
	log.Printf("Workflow done")

	// Wait for worker to exit
	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func MyWorkflow(ctx workflow.Context) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 100 * time.Second,
	})

	futs := make([]workflow.Future, 1000)
	for i := range futs {
		futs[i] = workflow.ExecuteActivity(ctx, MyActivity)
	}

	for i, fut := range futs {
		activityName := fmt.Sprintf("activity%v", i+1)
		if err := fut.Get(ctx, nil); err != nil {
			workflow.GetLogger(ctx).Error(activityName+" failed", "error", err)
		} else {
			workflow.GetLogger(ctx).Info(activityName + " complete")
		}
	}
	return nil
}

func fib(n uint) uint {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func MyActivity(ctx context.Context) error {
	then := time.Now()
	fib(43)
	fmt.Println("used: ", time.Since(then))
	return nil
}

func newPrometheusScope(c prometheus.Configuration) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		log.Fatalln("error creating prometheus reporter", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
		Prefix:          "cpu_intensive",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	log.Println("prometheus metrics scope created")
	return scope
}
