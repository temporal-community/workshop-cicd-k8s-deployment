package activities

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/temporal-community/workshop-cicd-k8s-deployment/shared"
	"go.temporal.io/sdk/activity"
)

// BuildDockerImage builds a Docker image from the specified context
func BuildDockerImage(ctx context.Context, req shared.DockerBuildRequest) (*shared.DockerBuildResponse, error) {
	logger := activity.GetLogger(ctx)
	startTime := time.Now()

	logger.Info("Starting Docker build",
		"image", req.ImageName,
		"tag", req.Tag,
		"context", req.BuildContext)

	// DEMO HELPER: Simulate random build failures
	if os.Getenv("SIMULATE_DOCKER_FAILURE") == "true" {
		if err := shared.SimulateFailure(0.3, "Docker daemon not responding"); err != nil {
			logger.Error("Simulated Docker build failure", "error", err)
			return nil, err
		}
	}

	// Construct the image tag
	imageTag := fmt.Sprintf("%s:%s", req.ImageName, req.Tag)

	// Build the Docker image
	cmd := exec.CommandContext(ctx, "docker", "buildx", "build",
		"-t", imageTag,
		"-f", req.Dockerfile,
		req.BuildContext)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Docker build failed",
			"error", err,
			"output", string(output))
		return nil, fmt.Errorf("docker build failed: %w\nOutput: %s", err, output)
	}

	// Get the image ID
	idCmd := exec.CommandContext(ctx, "docker", "images", "-q", imageTag)
	imageIDBytes, err := idCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get image ID: %w", err)
	}
	imageID := strings.TrimSpace(string(imageIDBytes))

	duration := time.Since(startTime)
	logger.Info("Docker build completed successfully",
		"imageID", imageID,
		"duration", duration)

	// Record heartbeat for long-running builds
	activity.RecordHeartbeat(ctx, "Build completed")

	return &shared.DockerBuildResponse{
		ImageID:   imageID,
		BuildTime: duration,
	}, nil
}

// TestDockerContainer runs tests against the built Docker image
func TestDockerContainer(ctx context.Context, req shared.DockerTestRequest) (*shared.DockerTestResponse, error) {
	logger := activity.GetLogger(ctx)
	startTime := time.Now()

	logger.Info("Starting Docker container tests",
		"image", req.ImageName,
		"tag", req.Tag)

	imageTag := fmt.Sprintf("%s:%s", req.ImageName, req.Tag)
	containerName := fmt.Sprintf("test-%s-%d", req.Tag, time.Now().Unix())

	// Start the container
	runCmd := exec.CommandContext(ctx, "docker", "run",
		"-d",
		"--name", containerName,
		"-p", "8080",
		imageTag)

	if output, err := runCmd.CombinedOutput(); err != nil {
		logger.Error("Failed to start test container",
			"error", err,
			"output", string(output))
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Ensure cleanup
	defer func() {
		stopCmd := exec.Command("docker", "rm", "-f", containerName)
		stopCmd.Run()
	}()

	// Get the mapped port
	portCmd := exec.CommandContext(ctx, "docker", "port", containerName, "8080")
	portOutput, err := portCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	// Extract port from output (format: "0.0.0.0:32768")
	portStr := strings.TrimSpace(string(portOutput))
	parts := strings.Split(portStr, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("unexpected port format: %s", portStr)
	}
	port := parts[len(parts)-1]

	// Wait for container to be ready
	time.Sleep(2 * time.Second)

	// Run integration tests
	testCmd := exec.CommandContext(ctx, "go", "test")
	testCmd.Dir = "sample-app"
	testCmd.Env = append(os.Environ(), fmt.Sprintf("BASE_URL=http://localhost:%s", port))

	testOutput, err := testCmd.CombinedOutput()
	passed := err == nil

	duration := time.Since(startTime)
	logger.Info("Docker tests completed",
		"passed", passed,
		"duration", duration,
		"output", string(testOutput))

	return &shared.DockerTestResponse{
		Passed:   passed,
		TestTime: duration,
		Output:   string(testOutput),
	}, nil
}

// PushToRegistry pushes the Docker image to the specified registry
func PushToRegistry(ctx context.Context, req shared.DockerPushRequest) (*shared.DockerPushResponse, error) {
	logger := activity.GetLogger(ctx)
	startTime := time.Now()

	logger.Info("Starting Docker push to registry",
		"image", req.ImageName,
		"tag", req.Tag,
		"registry", req.RegistryURL)

	// DEMO HELPER: Simulate occasional push failures
	if os.Getenv("SIMULATE_PUSH_FAILURE") == "true" {
		if err := shared.SimulateFailure(0.2, "Registry temporarily unavailable"); err != nil {
			logger.Error("Simulated registry push failure", "error", err)
			return nil, err
		}
	}

	localTag := fmt.Sprintf("%s:%s", req.ImageName, req.Tag)
	remoteTag := shared.FormatImageTag(req.RegistryURL, req.ImageName, req.Tag)

	// Tag the image for the remote registry
	tagCmd := exec.CommandContext(ctx, "docker", "tag", localTag, remoteTag)
	if output, err := tagCmd.CombinedOutput(); err != nil {
		logger.Error("Failed to tag image",
			"error", err,
			"output", string(output))
		return nil, fmt.Errorf("failed to tag image: %w", err)
	}

	// Push to registry
	pushCmd := exec.CommandContext(ctx, "docker", "push", remoteTag)
	pushOutput, err := pushCmd.CombinedOutput()
	if err != nil {
		logger.Error("Failed to push image",
			"error", err,
			"output", string(pushOutput))
		return nil, fmt.Errorf("failed to push image: %w\nOutput: %s", err, pushOutput)
	}

	// Extract digest from push output
	digest := extractDigest(string(pushOutput))

	duration := time.Since(startTime)
	logger.Info("Docker push completed successfully",
		"digest", digest,
		"duration", duration)

	// Record heartbeat for long pushes
	activity.RecordHeartbeat(ctx, "Push completed")

	return &shared.DockerPushResponse{
		Digest:   digest,
		PushTime: duration,
	}, nil
}

// extractDigest extracts the digest from docker push output
func extractDigest(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "digest:") && strings.Contains(line, "sha256:") {
			parts := strings.Split(line, "sha256:")
			if len(parts) >= 2 {
				return "sha256:" + strings.TrimSpace(parts[1])
			}
		}
	}
	return "unknown"
}
