package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	proxied "github.com/taonic/my-samples-go/proxied_activities"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "hello-proxy", worker.Options{})
	w.RegisterWorkflow(proxied.HelloProxyWorkflow)

	proxyActivity := proxied.NewProxyActivity(os.Getpid())
	w.RegisterActivity(proxyActivity.ProxyActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
