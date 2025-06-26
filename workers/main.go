package main

import (
	"log"
	"os"

	"github.com/temporal-workshop/cicd/activities"
	"github.com/temporal-workshop/cicd/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort: getTemporalHost(),
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	// Create worker
	w := worker.New(c, "cicd-task-queue", worker.Options{})

	// Register workflows
	w.RegisterWorkflow(workflows.BasicPipelineWorkflow)

	// Register activities
	w.RegisterActivity(activities.BuildDockerImage)
	w.RegisterActivity(activities.TestDockerContainer)
	w.RegisterActivity(activities.PushToRegistry)

	log.Println("Starting Temporal worker for Demo 1 - Basic Pipeline")
	log.Println("Worker listening on task queue: cicd-task-queue")

	// Start worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}
}

func getTemporalHost() string {
	host := os.Getenv("TEMPORAL_HOST")
	if host == "" {
		return "localhost:7233"
	}
	return host
}