// Shared types that mirror the Go implementation

export interface SendApprovalRequestRequest {
  Environment: string;
  ImageTag: string;
  StagingURL: string;
}

export interface SendApprovalRequestResponse {
  Success: boolean;
  NotificationID: string;
  Message: string;
}

export interface LogApprovalDecisionRequest {
  Approved: boolean;
  Approver: string;
  Reason: string;
  Timestamp: string;
}

export interface LogApprovalDecisionResponse {
  Success: boolean;
  Message: string;
}

export interface SendApprovalNotificationRequest {
  Approved: boolean;
  Environment: string;
  ImageTag: string;
  Approver: string;
  Reason: string;
}

export interface SendApprovalNotificationResponse {
  Success: boolean;
  Message: string;
}