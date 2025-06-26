package activities

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"

	"github.com/temporal-community/workshop-cicd-k8s-deployment/shared"
)

// KubernetesActivities provides Kubernetes deployment operations
type KubernetesActivities struct {
	Namespace string
}

// DeployToKubernetes deploys the application to Kubernetes
func (k *KubernetesActivities) DeployToKubernetes(ctx context.Context, req shared.DeployToKubernetesRequest) (*shared.DeployToKubernetesResponse, error) {
	logger := activity.GetLogger(ctx)
	info := activity.GetInfo(ctx)
	
	// Log activity start
	logger.Info("Starting Kubernetes deployment",
		"image", req.ImageTag,
		"environment", req.Environment,
		"namespace", k.getNamespace(req.Environment),
		"activityID", info.ActivityID,
		"attempt", info.Attempt)

	// Simulate deployment steps
	steps := []string{
		"Updating deployment manifest",
		"Applying deployment to cluster",
		"Waiting for pods to be ready",
		"Updating service configuration",
		"Verifying deployment health",
	}

	for i, step := range steps {
		// Check for cancellation
		if ctx.Err() != nil {
			return nil, temporal.NewApplicationError("deployment cancelled", "CANCELLED")
		}

		logger.Info(fmt.Sprintf("[%d/%d] %s", i+1, len(steps), step))
		
		// Simulate work
		time.Sleep(2 * time.Second)
		
		// Record heartbeat for long-running activity
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Completed step %d of %d", i+1, len(steps)))
	}

	// Generate deployment URL based on environment
	var deploymentURL string
	if req.Environment == "staging" {
		deploymentURL = fmt.Sprintf("http://staging.demo-app.local:8080")
	} else {
		deploymentURL = fmt.Sprintf("https://demo-app.production.com")
	}

	logger.Info("Kubernetes deployment completed successfully",
		"environment", req.Environment,
		"deploymentURL", deploymentURL)

	return &shared.DeployToKubernetesResponse{
		Success:       true,
		DeploymentURL: deploymentURL,
		Message:       fmt.Sprintf("Successfully deployed %s to %s", req.ImageTag, req.Environment),
		Timestamp:     time.Now(),
	}, nil
}

// CheckDeploymentStatus checks the status of a Kubernetes deployment
func (k *KubernetesActivities) CheckDeploymentStatus(ctx context.Context, req shared.CheckDeploymentStatusRequest) (*shared.CheckDeploymentStatusResponse, error) {
	logger := activity.GetLogger(ctx)
	
	logger.Info("Checking deployment status",
		"environment", req.Environment,
		"namespace", k.getNamespace(req.Environment))

	// Simulate status check
	time.Sleep(1 * time.Second)

	// In a real implementation, this would query the Kubernetes API
	return &shared.CheckDeploymentStatusResponse{
		Ready:        true,
		Replicas:     3,
		ReadyReplicas: 3,
		Message:      "All pods are running and ready",
	}, nil
}

// RollbackDeployment rolls back a Kubernetes deployment
func (k *KubernetesActivities) RollbackDeployment(ctx context.Context, req shared.RollbackDeploymentRequest) (*shared.RollbackDeploymentResponse, error) {
	logger := activity.GetLogger(ctx)
	
	logger.Info("Rolling back deployment",
		"environment", req.Environment,
		"reason", req.Reason)

	// Simulate rollback steps
	steps := []string{
		"Finding previous deployment revision",
		"Reverting deployment manifest",
		"Applying rollback to cluster",
		"Waiting for rollback to complete",
		"Verifying rollback success",
	}

	for i, step := range steps {
		logger.Info(fmt.Sprintf("[Rollback %d/%d] %s", i+1, len(steps), step))
		time.Sleep(1 * time.Second)
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Rollback step %d of %d", i+1, len(steps)))
	}

	logger.Info("Deployment rollback completed successfully")

	return &shared.RollbackDeploymentResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully rolled back %s deployment", req.Environment),
	}, nil
}

// GetServiceURL retrieves the service URL for a deployment
func (k *KubernetesActivities) GetServiceURL(ctx context.Context, req shared.GetServiceURLRequest) (*shared.GetServiceURLResponse, error) {
	logger := activity.GetLogger(ctx)
	
	logger.Info("Getting service URL",
		"environment", req.Environment,
		"serviceName", req.ServiceName)

	// Simulate service lookup
	time.Sleep(500 * time.Millisecond)

	// Generate URL based on environment
	var serviceURL string
	if req.Environment == "staging" {
		serviceURL = fmt.Sprintf("http://staging.%s.local:8080", req.ServiceName)
	} else {
		serviceURL = fmt.Sprintf("https://%s.production.com", req.ServiceName)
	}

	return &shared.GetServiceURLResponse{
		URL:     serviceURL,
		Ready:   true,
		Message: "Service is accessible",
	}, nil
}

// Helper method to get namespace based on environment
func (k *KubernetesActivities) getNamespace(environment string) string {
	if k.Namespace != "" {
		return k.Namespace
	}
	
	switch environment {
	case "staging":
		return "staging"
	case "production":
		return "production"
	default:
		return "default"
	}
}