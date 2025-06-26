package workflows

import (
	"fmt"
	"time"

	"github.com/temporal-community/workshop-cicd-k8s-deployment/activities"
	"github.com/temporal-community/workshop-cicd-k8s-deployment/shared"
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

// PipelineWithApprovalWorkflow extends the basic pipeline with Kubernetes deployment and human approval
func PipelineWithApprovalWorkflow(ctx workflow.Context, request shared.PipelineRequest) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting pipeline with approval workflow",
		"image", request.ImageName,
		"tag", request.Tag,
		"environment", request.Environment)

	// Configure activity options
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

	// Phase 1: Build, Test, and Push (reuse basic pipeline logic)
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

	// Phase 2: Deploy to Staging
	logger.Info("Phase 2: Deploying to staging environment")
	
	var k8sActivities *activities.KubernetesActivities
	deployReq := shared.DeployToKubernetesRequest{
		ImageTag:    fmt.Sprintf("%s:%s", request.ImageName, request.Tag),
		Environment: "staging",
	}

	var deployResp shared.DeployToKubernetesResponse
	err = workflow.ExecuteActivity(ctx, k8sActivities.DeployToKubernetes, deployReq).Get(ctx, &deployResp)
	if err != nil {
		logger.Error("Staging deployment failed", "error", err)
		return fmt.Errorf("staging deployment failed: %w", err)
	}

	logger.Info("Staging deployment successful",
		"url", deployResp.DeploymentURL)

	// Phase 3: Human Approval for Production
	if request.Environment == "production" {
		logger.Info("Phase 3: Requesting approval for production deployment")

		// Send approval request
		var approvalActivities *activities.ApprovalActivities
		approvalReq := shared.SendApprovalRequestRequest{
			Environment: "production",
			ImageTag:    fmt.Sprintf("%s:%s", request.ImageName, request.Tag),
			StagingURL:  deployResp.DeploymentURL,
		}

		var approvalResp shared.SendApprovalRequestResponse
		err = workflow.ExecuteActivity(ctx, approvalActivities.SendApprovalRequest, approvalReq).Get(ctx, &approvalResp)
		if err != nil {
			logger.Error("Failed to send approval request", "error", err)
			return fmt.Errorf("failed to send approval request: %w", err)
		}

		// Wait for approval signal
		logger.Info("Waiting for approval decision...")
		
		// Create a channel to receive the approval signal
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
		err = workflow.ExecuteActivity(ctx, approvalActivities.LogApprovalDecision, logReq).Get(ctx, &logResp)
		if err != nil {
			logger.Error("Failed to log approval decision", "error", err)
		}

		// Send notification about the decision
		notifyReq := shared.SendApprovalNotificationRequest{
			Approved:    approvalSignal.Approved,
			Environment: "production",
			ImageTag:    fmt.Sprintf("%s:%s", request.ImageName, request.Tag),
			Approver:    approvalSignal.Approver,
			Reason:      approvalSignal.Reason,
		}

		var notifyResp shared.SendApprovalNotificationResponse
		err = workflow.ExecuteActivity(ctx, approvalActivities.SendApprovalNotification, notifyReq).Get(ctx, &notifyResp)
		if err != nil {
			logger.Error("Failed to send approval notification", "error", err)
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
			ImageTag:    fmt.Sprintf("%s:%s", request.ImageName, request.Tag),
			Environment: "production",
		}

		var prodDeployResp shared.DeployToKubernetesResponse
		err = workflow.ExecuteActivity(ctx, k8sActivities.DeployToKubernetes, prodDeployReq).Get(ctx, &prodDeployResp)
		if err != nil {
			logger.Error("Production deployment failed", "error", err)
			return fmt.Errorf("production deployment failed: %w", err)
		}

		logger.Info("Production deployment successful",
			"url", prodDeployResp.DeploymentURL)
	}

	logger.Info("Pipeline with approval completed successfully",
		"totalDuration", workflow.Now(ctx).Sub(workflow.GetInfo(ctx).WorkflowStartTime))

	return nil
}
