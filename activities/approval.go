package activities

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"

	"github.com/temporal-community/workshop-cicd-k8s-deployment/shared"
)

// ApprovalActivities provides human approval operations
type ApprovalActivities struct{}

// SendApprovalRequest sends a notification for approval
func (a *ApprovalActivities) SendApprovalRequest(ctx context.Context, req shared.SendApprovalRequestRequest) (*shared.SendApprovalRequestResponse, error) {
	logger := activity.GetLogger(ctx)
	info := activity.GetInfo(ctx)

	logger.Info("Sending approval request",
		"environment", req.Environment,
		"imageTag", req.ImageTag,
		"stagingURL", req.StagingURL,
		"workflowID", info.WorkflowExecution.ID,
		"runID", info.WorkflowExecution.RunID)

	// In a real implementation, this would send notifications via Slack, email, etc.
	// For the demo, we'll just log the approval request details

	approvalMessage := fmt.Sprintf(`
==================================================
APPROVAL REQUIRED - Production Deployment
==================================================
Workflow ID: %s
Image Tag: %s
Environment: %s
Staging URL: %s

The application has been successfully deployed to staging.
Please review the staging deployment and approve or reject
the production deployment.

To approve:
  go run cmd/starter/main.go -action=approve -workflow=%s

To reject:
  go run cmd/starter/main.go -action=reject -workflow=%s

To check status:
  go run cmd/starter/main.go -action=status -workflow=%s
==================================================
`,
		info.WorkflowExecution.ID,
		req.ImageTag,
		req.Environment,
		req.StagingURL,
		info.WorkflowExecution.ID,
		info.WorkflowExecution.ID,
		info.WorkflowExecution.ID)

	logger.Info(approvalMessage)

	return &shared.SendApprovalRequestResponse{
		Success:        true,
		NotificationID: fmt.Sprintf("approval-%s-%d", info.WorkflowExecution.ID, time.Now().Unix()),
		Message:        "Approval request sent successfully",
	}, nil
}

// LogApprovalDecision logs the approval decision
func (a *ApprovalActivities) LogApprovalDecision(ctx context.Context, req shared.LogApprovalDecisionRequest) (*shared.LogApprovalDecisionResponse, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Logging approval decision",
		"approved", req.Approved,
		"approver", req.Approver,
		"reason", req.Reason,
		"timestamp", req.Timestamp)

	// In a real implementation, this might update a database, send notifications, etc.
	var message string
	if req.Approved {
		message = fmt.Sprintf("Deployment APPROVED by %s at %s", req.Approver, req.Timestamp.Format(time.RFC3339))
		if req.Reason != "" {
			message += fmt.Sprintf(" - Reason: %s", req.Reason)
		}
	} else {
		message = fmt.Sprintf("Deployment REJECTED by %s at %s", req.Approver, req.Timestamp.Format(time.RFC3339))
		if req.Reason != "" {
			message += fmt.Sprintf(" - Reason: %s", req.Reason)
		}
	}

	logger.Info(message)

	return &shared.LogApprovalDecisionResponse{
		Success: true,
		Message: message,
	}, nil
}

// SendApprovalNotification sends a notification about the approval decision
func (a *ApprovalActivities) SendApprovalNotification(ctx context.Context, req shared.SendApprovalNotificationRequest) (*shared.SendApprovalNotificationResponse, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Sending approval notification",
		"approved", req.Approved,
		"environment", req.Environment)

	// Build notification message
	var notificationMessage string
	if req.Approved {
		notificationMessage = fmt.Sprintf(`
==================================================
DEPLOYMENT APPROVED - Proceeding to Production
==================================================
Environment: %s
Image Tag: %s
Approved by: %s
Time: %s

The deployment has been approved and will now
proceed to the production environment.
==================================================
`,
			req.Environment,
			req.ImageTag,
			req.Approver,
			time.Now().Format(time.RFC3339))
	} else {
		notificationMessage = fmt.Sprintf(`
==================================================
DEPLOYMENT REJECTED - Workflow Cancelled
==================================================
Environment: %s
Image Tag: %s
Rejected by: %s
Reason: %s
Time: %s

The deployment has been rejected. The workflow
has been cancelled and no changes will be made
to the production environment.
==================================================
`,
			req.Environment,
			req.ImageTag,
			req.Approver,
			req.Reason,
			time.Now().Format(time.RFC3339))
	}

	logger.Info(notificationMessage)

	return &shared.SendApprovalNotificationResponse{
		Success: true,
		Message: "Notification sent successfully",
	}, nil
}
