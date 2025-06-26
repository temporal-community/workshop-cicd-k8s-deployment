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

// SimulateFailure randomly fails based on probability (for demos)
// DEMO HELPER - DO NOT USE IN PRODUCTION
func SimulateFailure(probability float32, errorMsg string) error {
	if rand.Float32() < probability {
		return fmt.Errorf("SIMULATED FAILURE: %s", errorMsg)
	}
	return nil
}
