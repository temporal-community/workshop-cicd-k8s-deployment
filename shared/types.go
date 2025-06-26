package shared

import "time"

// Pipeline request types
type PipelineRequest struct {
	ImageName    string
	Tag          string
	RegistryURL  string
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
	Passed   bool
	TestTime time.Duration
	Output   string
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