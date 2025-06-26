import { Context, log } from '@temporalio/activity';
import {
  SendApprovalRequestRequest,
  SendApprovalRequestResponse,
  LogApprovalDecisionRequest,
  LogApprovalDecisionResponse,
  SendApprovalNotificationRequest,
  SendApprovalNotificationResponse,
} from '../types';

/**
 * Sends a notification for approval
 * TypeScript implementation of the Go SendApprovalRequest activity
 */
export async function sendApprovalRequest(
  request: SendApprovalRequestRequest
): Promise<SendApprovalRequestResponse> {
  const logger = log;
  const info = Context.current().info;

  logger.info('Sending approval request', {
    environment: request.Environment,
    imageTag: request.ImageTag,
    stagingURL: request.StagingURL,
    workflowId: info.workflowExecution.workflowId,
    runId: info.workflowExecution.runId,
  });

  // In a real implementation, this would send notifications via Slack, email, etc.
  // For the demo, we'll just log the approval request details
  const approvalMessage = `
==================================================
APPROVAL REQUIRED - Production Deployment
==================================================
Workflow ID: ${info.workflowExecution.workflowId}
Image Tag: ${request.ImageTag}
Environment: ${request.Environment}
Staging URL: ${request.StagingURL}

The application has been successfully deployed to staging.
Please review the staging deployment and approve or reject
the production deployment.

To approve:
  go run cmd/starter/main.go -action=approve -workflow=${info.workflowExecution.workflowId}

To reject:
  go run cmd/starter/main.go -action=reject -workflow=${info.workflowExecution.workflowId}

To check status:
  go run cmd/starter/main.go -action=status -workflow=${info.workflowExecution.workflowId}
==================================================
`;

  logger.info(approvalMessage);

  return {
    Success: true,
    NotificationID: `approval-${info.workflowExecution.workflowId}-${Date.now()}`,
    Message: 'Approval request sent successfully',
  };
}

/**
 * Logs the approval decision
 * TypeScript implementation of the Go LogApprovalDecision activity
 */
export async function logApprovalDecision(
  request: LogApprovalDecisionRequest
): Promise<LogApprovalDecisionResponse> {
  const logger = log;

  logger.info('Logging approval decision', {
    approved: request.Approved,
    approver: request.Approver,
    reason: request.Reason,
    timestamp: request.Timestamp,
  });

  // In a real implementation, this might update a database, send notifications, etc.
  let message: string;
  if (request.Approved) {
    message = `Deployment APPROVED by ${request.Approver} at ${request.Timestamp}`;
    if (request.Reason) {
      message += ` - Reason: ${request.Reason}`;
    }
  } else {
    message = `Deployment REJECTED by ${request.Approver} at ${request.Timestamp}`;
    if (request.Reason) {
      message += ` - Reason: ${request.Reason}`;
    }
  }

  logger.info(message);

  return {
    Success: true,
    Message: message,
  };
}

/**
 * Sends a notification about the approval decision
 * TypeScript implementation of the Go SendApprovalNotification activity
 */
export async function sendApprovalNotification(
  request: SendApprovalNotificationRequest
): Promise<SendApprovalNotificationResponse> {
  const logger = log;

  logger.info('Sending approval notification', {
    approved: request.Approved,
    environment: request.Environment,
  });

  // Build notification message
  let notificationMessage: string;
  if (request.Approved) {
    notificationMessage = `
==================================================
DEPLOYMENT APPROVED - Proceeding to Production
==================================================
Environment: ${request.Environment}
Image Tag: ${request.ImageTag}
Approved by: ${request.Approver}
Time: ${new Date().toISOString()}

The deployment has been approved and will now
proceed to the production environment.
==================================================
`;
  } else {
    notificationMessage = `
==================================================
DEPLOYMENT REJECTED - Workflow Cancelled
==================================================
Environment: ${request.Environment}
Image Tag: ${request.ImageTag}
Rejected by: ${request.Approver}
Reason: ${request.Reason}
Time: ${new Date().toISOString()}

The deployment has been rejected. The workflow
has been cancelled and no changes will be made
to the production environment.
==================================================
`;
  }

  logger.info(notificationMessage);

  return {
    Success: true,
    Message: 'Notification sent successfully',
  };
}