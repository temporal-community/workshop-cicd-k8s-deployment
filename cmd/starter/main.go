package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/temporal-community/workshop-cicd-k8s-deployment/shared"
	"github.com/temporal-community/workshop-cicd-k8s-deployment/workflows"
	"go.temporal.io/sdk/client"
)

func main() {
	var (
		action       = flag.String("action", "create", "Action to perform: create, approve, reject, status")
		imageName    = flag.String("image", "demo-app", "Docker image name")
		tag          = flag.String("tag", "", "Docker image tag (defaults to v1.0.0)")
		registryURL  = flag.String("registry", "", "Container registry URL")
		environment  = flag.String("env", "staging", "Deployment environment: staging or production")
		buildContext = flag.String("context", "./sample-app", "Docker build context path")
		dockerfile   = flag.String("dockerfile", "Dockerfile", "Path to Dockerfile")
		workflowID   = flag.String("workflow", "", "Workflow ID for approval actions")
		approver     = flag.String("approver", "", "Name of the approver")
		reason       = flag.String("reason", "", "Reason for approval/rejection")
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

	switch *action {
	case "create":
		startPipeline(c, *imageName, *tag, *registryURL, *environment, *buildContext, *dockerfile)
	case "approve":
		if *workflowID == "" {
			log.Fatal("Workflow ID is required for approval action")
		}
		if *approver == "" {
			*approver = "demo-user"
		}
		sendApprovalSignal(c, *workflowID, true, *approver, *reason)
	case "reject":
		if *workflowID == "" {
			log.Fatal("Workflow ID is required for rejection action")
		}
		if *approver == "" {
			*approver = "demo-user"
		}
		if *reason == "" {
			*reason = "Deployment rejected"
		}
		sendApprovalSignal(c, *workflowID, false, *approver, *reason)
	case "status":
		if *workflowID == "" {
			log.Fatal("Workflow ID is required for status action")
		}
		getWorkflowStatus(c, *workflowID)
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

func startPipeline(c client.Client, imageName, tag, registryURL, environment, buildContext, dockerfile string) {
	// Create pipeline request
	request := shared.PipelineRequest{
		ImageName:    imageName,
		Tag:          tag,
		RegistryURL:  registryURL,
		Environment:  environment,
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

	// Use the unified CICDPipelineWorkflow with all features
	workflowFunc := workflows.CICDPipelineWorkflow

	// Start workflow
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflowFunc, request)
	if err != nil {
		log.Fatalf("Unable to start workflow: %v", err)
	}

	fmt.Printf("Started workflow:\n")
	fmt.Printf("  WorkflowID: %s\n", we.GetID())
	fmt.Printf("  RunID: %s\n", we.GetRunID())
	fmt.Printf("  Image: %s:%s\n", imageName, tag)
	fmt.Printf("  Registry: %s\n", registryURL)
	fmt.Printf("  Environment: %s\n", environment)
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

func sendApprovalSignal(c client.Client, workflowID string, approved bool, approver string, reason string) {
	// Create approval signal
	signal := shared.ApprovalSignal{
		Approved: approved,
		Approver: approver,
		Reason:   reason,
	}

	// Send signal to workflow
	err := c.SignalWorkflow(context.Background(), workflowID, "", "approval", signal)
	if err != nil {
		log.Fatalf("Unable to send approval signal: %v", err)
	}

	if approved {
		fmt.Printf("✅ Approval signal sent successfully!\n")
		fmt.Printf("  Workflow ID: %s\n", workflowID)
		fmt.Printf("  Approved by: %s\n", approver)
		if reason != "" {
			fmt.Printf("  Reason: %s\n", reason)
		}
	} else {
		fmt.Printf("❌ Rejection signal sent successfully!\n")
		fmt.Printf("  Workflow ID: %s\n", workflowID)
		fmt.Printf("  Rejected by: %s\n", approver)
		fmt.Printf("  Reason: %s\n", reason)
	}
}


func getWorkflowStatus(c client.Client, workflowID string) {
	// Get workflow description
	resp, err := c.DescribeWorkflowExecution(context.Background(), workflowID, "")
	if err != nil {
		log.Fatalf("Unable to get workflow status: %v", err)
	}

	fmt.Printf("Workflow Status:\n")
	fmt.Printf("  Workflow ID: %s\n", workflowID)
	fmt.Printf("  Status: %s\n", resp.WorkflowExecutionInfo.Status)
	fmt.Printf("  Start Time: %s\n", resp.WorkflowExecutionInfo.StartTime)
	
	if resp.WorkflowExecutionInfo.CloseTime != nil {
		fmt.Printf("  Close Time: %s\n", resp.WorkflowExecutionInfo.CloseTime)
	} else {
		fmt.Printf("  Running for: %s\n", time.Since(resp.WorkflowExecutionInfo.StartTime.AsTime()))
	}

	// Show pending activities if any
	if len(resp.PendingActivities) > 0 {
		fmt.Printf("\nPending Activities:\n")
		for _, activity := range resp.PendingActivities {
			fmt.Printf("  - %s (attempts: %d)\n", activity.ActivityType.Name, activity.Attempt)
		}
	}
}