# Demo 1: Basic Docker Pipeline

This demo showcases a basic Temporal workflow that builds, tests, and pushes Docker images with built-in retry policies.

## Overview

The workflow demonstrates:
- Sequential activity execution
- Retry policies for transient failures
- Comprehensive logging and visibility
- Idempotent operations

## Architecture

```
BasicPipelineWorkflow
├── BuildDockerImage (Activity)
├── TestDockerContainer (Activity)
└── PushToRegistry (Activity)
```

## Running the Demo

### Prerequisites

1. Temporal server running:
   ```bash
   temporal server start-dev
   ```

2. Docker daemon running
3. Registry credentials configured (if using remote registry)

### Step 1: Start the Worker

```bash
go run workers/main.go
```

You should see:
```
Starting Temporal worker for Demo 1 - Basic Pipeline
Worker listening on task queue: cicd-task-queue
```

### Step 2: Trigger the Workflow

Basic usage:
```bash
go run cmd/starter/main.go -image=demo-app -tag=v1.0.0 -dockerfile=sample-app/Dockerfile -registry=registry.digitalocean.com/ziggys-container
```

### Step 3: View in Temporal UI

Open http://localhost:8233 and navigate to the workflow.

## Demonstrating Features

### 1. Retry Policies

Simulate Docker build failures:
```bash
export SIMULATE_DOCKER_FAILURE=true
go run cmd/starter/main.go -image=demo-app -tag=v1.0.1
```

Watch the Temporal UI to see:
- Activity retries with exponential backoff
- Retry attempts in the event history
- Eventual success after transient failures

### 2. Activity Timeouts

Each activity has a 10-minute timeout. If Docker operations take longer, the activity will timeout and retry.

### 3. Workflow Visibility

In the Temporal UI, demonstrate:
- Real-time workflow progress
- Activity inputs and outputs
- Event history showing each step
- Workflow execution time

## Key Code Points

### Workflow Implementation (workflows/pipeline.go)

- Sequential execution pattern
- Activity options with retry configuration
- Structured logging for visibility

### Activity Implementation (activities/docker.go)

- Real Docker operations (build, test, push)
- Error simulation for demos
- Heartbeat recording for long operations
- Idempotent design

### Worker Registration (workers/main.go)

- Simple worker setup
- Activity registration
- Task queue configuration

## Common Issues

### Docker daemon not accessible
```bash
# Verify Docker is running
docker ps
```

### Registry authentication issues
```bash
# Login to registry
docker login $REGISTRY_URL
```

### Port conflicts during testing
The test activity dynamically allocates ports to avoid conflicts.

## Next Steps

After completing this demo, proceed to Demo 2 where we'll add:
- Kubernetes deployment
- Human approval workflows
- Signal handling