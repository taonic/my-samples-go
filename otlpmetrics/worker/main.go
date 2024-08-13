package main

import (
	"context"
	"log"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"

	metrics "github.com/taonic/my-samples-go/otlpmetrics"
)

func main() {
	ctx := context.Background()
	var err error
	useGRPC := true
	var exp metric.Exporter
	if useGRPC {
		conn, err := grpc.Dial("localhost:4317", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		exp, err = otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
		if err != nil {
			panic(err)
		}
	} else {
		exp, err = otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL("http://localhost:4318"))
		if err != nil {
			panic(err)
		}
	}
	meterProvider := metric.NewMeterProvider(metric.WithReader(
		metric.NewPeriodicReader(exp, metric.WithInterval(10*time.Second)),
	))
	c, err := client.Dial(client.Options{
		MetricsHandler: opentelemetry.NewMetricsHandler(
			opentelemetry.MetricsHandlerOptions{
				Meter: meterProvider.Meter("temporal-sdk-go"),
			},
		),
	})

	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "metrics", worker.Options{})

	w.RegisterWorkflow(metrics.Workflow)
	w.RegisterActivity(metrics.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
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
		Prefix:          "temporal",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	log.Println("prometheus metrics scope created")
	return scope
}
