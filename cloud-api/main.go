package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/taonic/my-samples-go/cloud-api/client/api"
	"github.com/taonic/my-samples-go/cloud-api/protogen/temporal/api/cloud/cloudservice/v1"
)

const (
	temporalCloudAPIAddress    = "saas-api.tmprl.cloud:443"
	temporalCloudAPIKeyEnvName = "TEMPORAL_CLOUD_API_KEY"
)

var namespace string

func init() {
	flag.StringVar(&namespace, "namespace", "", "Temporal namespace")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	apikey := os.Getenv(temporalCloudAPIKeyEnvName)
	if apikey == "" {
		return fmt.Errorf("apikey not provided, set environment variable '%s' with apikey you want to use", temporalCloudAPIKeyEnvName)
	}

	flag.Parse()

	if namespace == "" {
		return fmt.Errorf("--namespace flag is required")
	}

	conn, err := api.NewConnectionWithAPIKey(temporalCloudAPIAddress, false, apikey)
	if err != nil {
		return fmt.Errorf("failed to create cloud api connection: %+v", err)
	}

	input := cloudservice.GetNamespaceRequest{
		Namespace: namespace,
	}
	result, err := cloudservice.NewCloudServiceClient(conn).GetNamespace(context.Background(), &input)
	if err != nil {
		return err
	}
	log.Printf("Custom search attributes: %+v", result.Namespace.Spec.CustomSearchAttributes)

	return nil
}
