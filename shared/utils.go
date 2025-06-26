package shared

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateWorkflowID creates a unique workflow ID
func GenerateWorkflowID(prefix string) string {
	timestamp := time.Now().Unix()
	random := rand.Intn(9999)
	return fmt.Sprintf("%s-%d-%04d", prefix, timestamp, random)
}

// FormatImageTag formats the full image name with registry
func FormatImageTag(registry, image, tag string) string {
	if registry == "" {
		return fmt.Sprintf("%s:%s", image, tag)
	}
	return fmt.Sprintf("%s/%s:%s", registry, image, tag)
}

// IsProductionEnvironment checks if the environment is production
func IsProductionEnvironment(env string) bool {
	return env == "production" || env == "prod"
}

// IsStagingEnvironment checks if the environment is staging
func IsStagingEnvironment(env string) bool {
	return env == "staging" || env == "stage"
}

// GetNamespaceForEnvironment returns the k8s namespace for an environment
func GetNamespaceForEnvironment(env string) string {
	if IsProductionEnvironment(env) {
		return "production"
	}
	if IsStagingEnvironment(env) {
		return "staging"
	}
	return "default"
}

// SimulateFailure randomly fails based on probability (for demos)
// DEMO HELPER - DO NOT USE IN PRODUCTION
func SimulateFailure(probability float32, errorMsg string) error {
	if rand.Float32() < probability {
		return fmt.Errorf("SIMULATED FAILURE: %s", errorMsg)
	}
	return nil
}

// IsWithinDeploymentWindow checks if current time is within deployment window
func IsWithinDeploymentWindow(startHour, endHour int) bool {
	now := time.Now()
	currentHour := now.Hour()
	
	// Simple logic for demo - in production this would be more sophisticated
	if startHour <= endHour {
		return currentHour >= startHour && currentHour < endHour
	}
	// Handle overnight windows (e.g., 22:00 - 02:00)
	return currentHour >= startHour || currentHour < endHour
}

// GetDeploymentWindowWaitTime calculates how long to wait for next window
func GetDeploymentWindowWaitTime(startHour int) time.Duration {
	now := time.Now()
	currentHour := now.Hour()
	
	waitHours := startHour - currentHour
	if waitHours <= 0 {
		waitHours += 24
	}
	
	// For demo purposes, return seconds instead of hours
	// In production, this would return actual hours
	return time.Duration(waitHours) * time.Second
}