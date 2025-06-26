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

	// Build the Docker image for local testing (current platform only)
	// Multi-arch build will happen during push phase
	logger.Info("Building image for local testing")
	
	cmd := exec.CommandContext(ctx, "docker", "buildx", "build",
		"-t", imageTag,
		"-f", req.Dockerfile,
		"--load", // Load for local testing
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

	remoteTag := shared.FormatImageTag(req.RegistryURL, req.ImageName, req.Tag)

	// Build and push multi-architecture image directly to registry
	logger.Info("Building and pushing multi-architecture image", 
		"platforms", "linux/amd64,linux/arm64",
		"remoteTag", remoteTag)

	// Ensure we have a multi-platform capable builder
	builderName := "multiarch-builder"
	
	// Remove existing builder first to avoid conflicts
	removeBuilderCmd := exec.CommandContext(ctx, "docker", "buildx", "rm", builderName)
	removeBuilderCmd.Run() // Ignore errors if builder doesn't exist
	
	// Create fresh builder
	createBuilderCmd := exec.CommandContext(ctx, "docker", "buildx", "create", 
		"--name", builderName, 
		"--driver", "docker-container",
		"--use")
	createOutput, createErr := createBuilderCmd.CombinedOutput()
	
	if createErr != nil {
		logger.Error("Failed to create multi-arch builder", 
			"error", createErr,
			"output", string(createOutput))
		return nil, fmt.Errorf("failed to create multi-arch builder: %w", createErr)
	}
	
	logger.Info("Created fresh multi-arch builder", "name", builderName)

	// Rebuild for multi-architecture and push directly to registry
	// Use a unique tag to avoid conflicts with existing untagged images
	timestampTag := fmt.Sprintf("%s-%d", req.Tag, time.Now().Unix())
	multiArchTag := shared.FormatImageTag(req.RegistryURL, req.ImageName, timestampTag)
	
	buildCmd := exec.CommandContext(ctx, "docker", "buildx", "build",
		"--platform", "linux/amd64,linux/arm64",
		"-t", multiArchTag,
		"-f", req.Dockerfile,
		"--no-cache", // Force clean build to ensure CGO_ENABLED=0 fix is applied
		"--push", // Push directly to registry
		req.BuildContext)

	pushOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		logger.Warn("Multi-arch build failed, falling back to single-arch build and push",
			"error", err,
			"output", string(pushOutput))
		
		// Fallback: Tag and push single-arch image
		localTag := fmt.Sprintf("%s:%s", req.ImageName, req.Tag)
		
		// Tag the image for the remote registry
		tagCmd := exec.CommandContext(ctx, "docker", "tag", localTag, remoteTag)
		if tagOutput, tagErr := tagCmd.CombinedOutput(); tagErr != nil {
			logger.Error("Failed to tag image for fallback",
				"error", tagErr,
				"output", string(tagOutput))
			return nil, fmt.Errorf("failed to tag image: %w", tagErr)
		}

		// Push single-arch image
		fallbackPushCmd := exec.CommandContext(ctx, "docker", "push", remoteTag)
		fallbackOutput, fallbackErr := fallbackPushCmd.CombinedOutput()
		if fallbackErr != nil {
			logger.Error("Fallback push also failed",
				"error", fallbackErr,
				"output", string(fallbackOutput))
			return nil, fmt.Errorf("failed to push image (both multi-arch and fallback failed): %w\nOutput: %s", fallbackErr, fallbackOutput)
		}
		
		logger.Info("Successfully pushed single-arch image as fallback")
		pushOutput = fallbackOutput
	} else {
		// Multi-arch build succeeded, now tag it with the original tag
		logger.Info("Multi-arch build succeeded, creating additional tag", "originalTag", remoteTag)
		
		// Use buildx imagetools to create an additional tag pointing to the same manifest
		tagCmd := exec.CommandContext(ctx, "docker", "buildx", "imagetools", "create", 
			"-t", remoteTag, 
			multiArchTag)
		if tagOutput, tagErr := tagCmd.CombinedOutput(); tagErr != nil {
			logger.Warn("Failed to create additional tag, but multi-arch push succeeded",
				"error", tagErr,
				"output", string(tagOutput))
		} else {
			logger.Info("Successfully created additional tag", "tag", remoteTag)
		}
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
