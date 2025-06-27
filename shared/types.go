package shared

import "time"

// Pipeline request types
type PipelineRequest struct {
	ImageName    string
	Tag          string
	RegistryURL  string
	Environment  string // staging or production
	BuildContext string
	Dockerfile   string
}

// Docker activity types
type DockerBuildRequest struct {
	ImageName    string
	Tag          string
	BuildContext string
	Dockerfile   string
}

type DockerBuildResponse struct {
	ImageID   string
	BuildTime time.Duration
}

type DockerTestRequest struct {
	ImageName string
	Tag       string
}

type DockerTestResponse struct {
	Passed    bool
	TestTime  time.Duration
	Output    string
}

type DockerPushRequest struct {
	ImageName    string
	Tag          string
	RegistryURL  string
	BuildContext string
	Dockerfile   string
}

type DockerPushResponse struct {
	Digest   string
	PushTime time.Duration
}

// Kubernetes activity types
type KubernetesDeployRequest struct {
	ImageName   string
	Tag         string
	RegistryURL string
	Namespace   string
	Replicas    int32
}

type KubernetesDeployResponse struct {
	DeploymentName string
	ServiceURL     string
	Status         string
}

// Approval types
type ApprovalRequest struct {
	WorkflowID   string
	RunID        string
	Environment  string
	ImageName    string
	Tag          string
	StagingURL   string
	RequestedBy  string
	RequestedAt  time.Time
}

type ApprovalResponse struct {
	Approved    bool
	ApprovedBy  string
	ApprovedAt  time.Time
	Comments    string
}

// Monitoring types
type HealthCheckRequest struct {
	ServiceURL string
	Timeout    time.Duration
}

type HealthCheckResponse struct {
	Healthy       bool
	ResponseTime  time.Duration
	StatusCode    int
	Error         string
}

type RollbackRequest struct {
	DeploymentName string
	Namespace      string
	PreviousTag    string
}

// Workflow states
type WorkflowState struct {
	Status           string
	CurrentActivity  string
	DockerBuildDone  bool
	DockerTestDone   bool
	DockerPushDone   bool
	StagingDeployed  bool
	ApprovalReceived bool
	ProductionReady  bool
}

// Additional Kubernetes activity types
type DeployToKubernetesRequest struct {
	ImageTag    string
	Environment string // staging or production
}

type DeployToKubernetesResponse struct {
	Success       bool
	DeploymentURL string
	Message       string
	Timestamp     time.Time
}

type CheckDeploymentStatusRequest struct {
	Environment string
}

type CheckDeploymentStatusResponse struct {
	Ready         bool
	Replicas      int32
	ReadyReplicas int32
	Message       string
}


type GetServiceURLRequest struct {
	Environment string
	ServiceName string
}

type GetServiceURLResponse struct {
	URL     string
	Ready   bool
	Message string
}

// Additional Approval activity types
type SendApprovalRequestRequest struct {
	Environment string
	ImageTag    string
	StagingURL  string
}

type SendApprovalRequestResponse struct {
	Success        bool
	NotificationID string
	Message        string
}

type LogApprovalDecisionRequest struct {
	Approved  bool
	Approver  string
	Reason    string
	Timestamp time.Time
}

type LogApprovalDecisionResponse struct {
	Success bool
	Message string
}

type SendApprovalNotificationRequest struct {
	Approved    bool
	Environment string
	ImageTag    string
	Approver    string
	Reason      string
}

type SendApprovalNotificationResponse struct {
	Success bool
	Message string
}

// Approval signal types
type ApprovalSignal struct {
	Approved bool
	Approver string
	Reason   string
}


