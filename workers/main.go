package main

import (
	"log"
	"os"

	"github.com/temporal-community/workshop-cicd-k8s-deployment/activities"
	"github.com/temporal-community/workshop-cicd-k8s-deployment/workflows"
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
	w := worker.New(c, "cicd-task-queue-go", worker.Options{})

	// Register workflows
	w.RegisterWorkflow(workflows.CICDPipelineWorkflow)

	// Register Docker activities
	w.RegisterActivity(activities.BuildDockerImage)
	w.RegisterActivity(activities.TestDockerContainer)
	w.RegisterActivity(activities.PushToRegistry)

	// Kubernetes activities handled by Python worker

	// Approval activities handled by TypeScript worker


	log.Println("Starting Temporal worker for CI/CD Pipeline (Polyglot Mode)")
	log.Println("Worker listening on task queue: cicd-task-queue-go")
	log.Println("Registered workflows:")
	log.Println("  - CICDPipelineWorkflow (unified workflow with all features)")
	log.Println("Registered activities:")
	log.Println("  - Docker: Build, Test, Push")
	log.Println("  - Kubernetes: Handled by Python worker")
	log.Println("  - Approval: Handled by TypeScript worker")

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