package workflows

import (
	"fmt"
	"time"

	"github.com/temporal-workshop/cicd/activities"
	"github.com/temporal-workshop/cicd/shared"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// BasicPipelineWorkflow implements a simple sequential CI/CD pipeline
func BasicPipelineWorkflow(ctx workflow.Context, request shared.PipelineRequest) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting basic pipeline workflow", 
		"image", request.ImageName,
		"tag", request.Tag,
		"registry", request.RegistryURL)

	// Configure activity options with retry policy
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Step 1: Build Docker image
	logger.Info("Starting Docker build")
	var buildResp shared.DockerBuildResponse
	buildReq := shared.DockerBuildRequest{
		ImageName:    request.ImageName,
		Tag:          request.Tag,
		BuildContext: request.BuildContext,
		Dockerfile:   "Dockerfile",
	}
	
	err := workflow.ExecuteActivity(ctx, activities.BuildDockerImage, buildReq).Get(ctx, &buildResp)
	if err != nil {
		logger.Error("Docker build failed", "error", err)
		return fmt.Errorf("docker build failed: %w", err)
	}
	logger.Info("Docker build completed", 
		"imageID", buildResp.ImageID,
		"duration", buildResp.BuildTime)

	// Step 2: Test Docker container
	logger.Info("Starting Docker tests")
	var testResp shared.DockerTestResponse
	testReq := shared.DockerTestRequest{
		ImageName: request.ImageName,
		Tag:       request.Tag,
	}
	
	err = workflow.ExecuteActivity(ctx, activities.TestDockerContainer, testReq).Get(ctx, &testResp)
	if err != nil {
		logger.Error("Docker tests failed", "error", err)
		return fmt.Errorf("docker tests failed: %w", err)
	}
	logger.Info("Docker tests completed", 
		"passed", testResp.Passed,
		"duration", testResp.TestTime)

	if !testResp.Passed {
		return fmt.Errorf("docker tests failed: %s", testResp.Output)
	}

	// Step 3: Push to registry
	logger.Info("Starting Docker push")
	var pushResp shared.DockerPushResponse
	pushReq := shared.DockerPushRequest{
		ImageName:   request.ImageName,
		Tag:         request.Tag,
		RegistryURL: request.RegistryURL,
	}
	
	err = workflow.ExecuteActivity(ctx, activities.PushToRegistry, pushReq).Get(ctx, &pushResp)
	if err != nil {
		logger.Error("Docker push failed", "error", err)
		return fmt.Errorf("docker push failed: %w", err)
	}
	logger.Info("Docker push completed", 
		"digest", pushResp.Digest,
		"duration", pushResp.PushTime)

	logger.Info("Pipeline completed successfully", 
		"totalDuration", workflow.Now(ctx).Sub(workflow.GetInfo(ctx).WorkflowStartTime))
	
	return nil
}