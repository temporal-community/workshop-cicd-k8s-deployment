package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var baseURL string

func init() {
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
}

func TestHealthEndpoint(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if strings.TrimSpace(string(body)) != "OK" {
		t.Errorf("Expected 'OK', got '%s'", strings.TrimSpace(string(body)))
	}
}

func TestRootEndpoint(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(baseURL)
	if err != nil {
		t.Fatalf("Root endpoint test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), "Hello World") {
		t.Errorf("Expected 'Hello World' in response, got: %s", string(body))
	}
}
