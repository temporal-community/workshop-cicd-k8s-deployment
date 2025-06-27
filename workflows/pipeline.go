package workflows

import (
	"fmt"
	"time"

	"github.com/temporal-community/workshop-cicd-k8s-deployment/activities"
	"github.com/temporal-community/workshop-cicd-k8s-deployment/shared"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// CICDPipelineWorkflow implements a comprehensive CI/CD pipeline that adapts based on environment
func CICDPipelineWorkflow(ctx workflow.Context, request shared.PipelineRequest) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting CI/CD pipeline workflow",
		"image", request.ImageName,
		"tag", request.Tag,
		"registry", request.RegistryURL,
		"environment", request.Environment)

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

	// Phase 1: Build, Test, and Push (always happens)
	logger.Info("Phase 1: Docker build, test, and push")

	// Step 1: Build Docker image
	logger.Info("Starting Docker build")
	var buildResp shared.DockerBuildResponse
	buildReq := shared.DockerBuildRequest{
		ImageName:    request.ImageName,
		Tag:          request.Tag,
		BuildContext: request.BuildContext,
		Dockerfile:   request.Dockerfile,
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
		ImageName:    request.ImageName,
		Tag:          request.Tag,
		RegistryURL:  request.RegistryURL,
		BuildContext: request.BuildContext,
		Dockerfile:   request.Dockerfile,
	}

	err = workflow.ExecuteActivity(ctx, activities.PushToRegistry, pushReq).Get(ctx, &pushResp)
	if err != nil {
		logger.Error("Docker push failed", "error", err)
		return fmt.Errorf("docker push failed: %w", err)
	}
	logger.Info("Docker push completed",
		"digest", pushResp.Digest,
		"duration", pushResp.PushTime)

	// Construct full image path with registry
	var fullImagePath string
	if request.RegistryURL != "" {
		fullImagePath = fmt.Sprintf("%s/%s:%s", request.RegistryURL, request.ImageName, request.Tag)
	} else {
		fullImagePath = fmt.Sprintf("%s:%s", request.ImageName, request.Tag)
	}

	// Phase 2: Deploy to Staging (happens for staging and production environments)
	if request.Environment == "staging" || request.Environment == "production" {
		logger.Info("Phase 2: Deploying to staging environment")

		deployReq := shared.DeployToKubernetesRequest{
			ImageTag:    fullImagePath,
			Environment: "staging",
		}

		var deployResp shared.DeployToKubernetesResponse
		err = workflow.ExecuteActivity(ctx, "DeployToKubernetes", deployReq).Get(ctx, &deployResp)
		if err != nil {
			logger.Error("Staging deployment failed", "error", err)
			return fmt.Errorf("staging deployment failed: %w", err)
		}

		logger.Info("Staging deployment successful", "url", deployResp.DeploymentURL)

		// Phase 3: Production deployment with approval (if production environment)
		if request.Environment == "production" {
			err = deployToProduction(ctx, logger, fullImagePath, deployResp.DeploymentURL)
			if err != nil {
				return err
			}
		}
	}

	logger.Info("CI/CD pipeline completed successfully",
		"totalDuration", workflow.Now(ctx).Sub(workflow.GetInfo(ctx).WorkflowStartTime))

	return nil
}

// deployToProduction handles the production deployment with approval
func deployToProduction(ctx workflow.Context, logger log.Logger, fullImagePath, stagingURL string) error {
	logger.Info("Phase 3: Requesting approval for production deployment")

	// Send approval request
	approvalReq := shared.SendApprovalRequestRequest{
		Environment: "production",
		ImageTag:    fullImagePath,
		StagingURL:  stagingURL,
	}

	var approvalResp shared.SendApprovalRequestResponse
	err := workflow.ExecuteActivity(ctx, "SendApprovalRequest", approvalReq).Get(ctx, &approvalResp)
	if err != nil {
		logger.Error("Failed to send approval request", "error", err)
		return fmt.Errorf("failed to send approval request: %w", err)
	}

	// Wait for approval signal
	logger.Info("Waiting for approval decision...")
	approvalChannel := workflow.GetSignalChannel(ctx, "approval")
	var approvalSignal shared.ApprovalSignal
	approvalChannel.Receive(ctx, &approvalSignal)

	// Log the approval decision
	logReq := shared.LogApprovalDecisionRequest{
		Approved:  approvalSignal.Approved,
		Approver:  approvalSignal.Approver,
		Reason:    approvalSignal.Reason,
		Timestamp: workflow.Now(ctx),
	}

	var logResp shared.LogApprovalDecisionResponse
	err = workflow.ExecuteActivity(ctx, "LogApprovalDecision", logReq).Get(ctx, &logResp)
	if err != nil {
		logger.Error("Failed to log approval decision", "error", err)
	}

	// Check if approved
	if !approvalSignal.Approved {
		logger.Info("Production deployment rejected",
			"rejectedBy", approvalSignal.Approver,
			"reason", approvalSignal.Reason)
		return fmt.Errorf("production deployment rejected by %s: %s",
			approvalSignal.Approver, approvalSignal.Reason)
	}

	// Phase 4: Deploy to Production
	logger.Info("Phase 4: Deploying to production environment")

	prodDeployReq := shared.DeployToKubernetesRequest{
		ImageTag:    fullImagePath,
		Environment: "production",
	}

	var prodDeployResp shared.DeployToKubernetesResponse
	err = workflow.ExecuteActivity(ctx, "DeployToKubernetes", prodDeployReq).Get(ctx, &prodDeployResp)
	if err != nil {
		logger.Error("Production deployment failed", "error", err)
		return fmt.Errorf("production deployment failed: %w", err)
	}

	logger.Info("Production deployment successful", "url", prodDeployResp.DeploymentURL)

	return nil
}
