package shared

import "time"

// Pipeline request types
type PipelineRequest struct {
	ImageName    string
	Tag          string
	RegistryURL  string
	Environment  string // staging or production
	BuildContext string
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
	ImageName   string
	Tag         string
	RegistryURL string
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