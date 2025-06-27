package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/temporal-community/workshop-cicd-k8s-deployment/shared"
	"github.com/temporal-community/workshop-cicd-k8s-deployment/workflows"
	"go.temporal.io/sdk/client"
)

func main() {
	var (
		imageName    = flag.String("image", "demo-app", "Docker image name")
		tag          = flag.String("tag", "", "Docker image tag (defaults to v1.0.0)")
		registryURL  = flag.String("registry", "", "Container registry URL")
		buildContext = flag.String("context", "./sample-app", "Docker build context path")
		dockerfile   = flag.String("dockerfile", "Dockerfile", "Path to Dockerfile")
	)
	flag.Parse()

	// Set defaults
	if *tag == "" {
		*tag = "v1.0.0"
	}

	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort: getTemporalHost(),
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	// Start the pipeline
	startPipeline(c, *imageName, *tag, *registryURL, *buildContext, *dockerfile)
}

func startPipeline(c client.Client, imageName, tag, registryURL, buildContext, dockerfile string) {
	// Create pipeline request
	request := shared.PipelineRequest{
		ImageName:    imageName,
		Tag:          tag,
		RegistryURL:  registryURL,
		BuildContext: buildContext,
		Dockerfile:   dockerfile,
	}

	// Generate workflow ID
	workflowID := shared.GenerateWorkflowID("pipeline")

	// Workflow options
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "cicd-task-queue",
	}

	// Start workflow
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.CICDPipelineWorkflow, request)
	if err != nil {
		log.Fatalf("Unable to start workflow: %v", err)
	}

	fmt.Printf("Started workflow:\n")
	fmt.Printf("  WorkflowID: %s\n", we.GetID())
	fmt.Printf("  RunID: %s\n", we.GetRunID())
	fmt.Printf("  Image: %s:%s\n", imageName, tag)
	fmt.Printf("  Registry: %s\n", registryURL)
	fmt.Printf("\n")
	fmt.Printf("View in Temporal UI: http://localhost:8233/namespaces/default/workflows/%s\n", we.GetID())
}

func getTemporalHost() string {
	host := os.Getenv("TEMPORAL_HOST")
	if host == "" {
		return "localhost:7233"
	}
	return host
}