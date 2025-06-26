# Temporal CI/CD Workshop - Part 1: Docker Pipeline

A hands-on demonstration of how Temporal enables reliable CI/CD pipelines through a simple Docker build, test, and push workflow.

## What This Demo Shows

This Part 1 implementation demonstrates a **CI/CD pipeline** that showcases core Temporal concepts:

1. **Sequential workflow execution**: Docker build → test → push
2. **Activity retry policies**: Automatic retries with exponential backoff
3. **Workflow visibility**: Complete observability through Temporal UI
4. **Error handling**: Graceful failure handling and recovery

## Key Concepts Demonstrated

- **Workflows as Code**: Define your pipeline logic in Go
- **Activities**: Discrete units of work (build, test, push)
- **Retry Policies**: Automatic retry with configurable backoff
- **Deterministic Execution**: Workflows that can be replayed exactly
- **Temporal UI**: Visual workflow execution tracking

## Prerequisites

- **Docker Desktop** installed and running
- **Go 1.21+**
- **Temporal CLI** (`brew install temporal` or download from https://temporal.io/downloads)

## Quick Start

### 1. Start Temporal Server
```bash
temporal server start-dev
```

This starts a local Temporal server with Web UI at http://localhost:8233

### 2. Start the Worker
In a new terminal:
```bash
go run workers/main.go
```

You should see:
```
Starting Temporal worker for CI/CD Pipeline
Worker listening on task queue: cicd-task-queue
Registered workflows:
  - CICDPipelineWorkflow
Registered activities:
  - Docker: Build, Test, Push
```

### 3. Run the Pipeline
In another terminal:
```bash
go run cmd/starter/main.go \
  -image=demo-app \
  -tag=v1.0.0 \
  -registry=myregistry.io
```

Output:
```
Started workflow:
  WorkflowID: pipeline-1234567890-0001
  RunID: <unique-run-id>
  Image: demo-app:v1.0.0
  Registry: myregistry.io

View in Temporal UI: http://localhost:8233/namespaces/default/workflows/pipeline-1234567890-0001
```

### 4. Monitor Progress

1. **Open Temporal UI**: http://localhost:8233
2. **Click on your workflow ID**
3. **Watch the execution**:
   - ✓ BuildDockerImage activity
   - ✓ TestDockerContainer activity  
   - ✓ PushToRegistry activity
   - ✓ Workflow completed

## Demo

### Run Worker
```bash
# Run with default settings
go run cmd/starter/main.go
```

### Start Workflow
```bash
 go run cmd/starter/main.go \
  -image=demo-app \
  -tag=2.0.0 \
  -dockerfile=sample-app/Dockerfile \
  -registry=REGISTRY_URL
```



## Available Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-image` | Docker image name | `demo-app` |
| `-tag` | Docker image tag | `v1.0.0` |
| `-registry` | Container registry URL | (empty) |
| `-dockerfile` | Path to Dockerfile | `Dockerfile` |


## Project Structure

```
├── README.md                    # This file
├── activities/                  # Activity implementations
│   └── docker.go               # Docker build/test/push operations
├── cmd/starter/main.go         # CLI to start workflows
├── workers/main.go             # Temporal worker
├── workflows/pipeline.go       # Pipeline workflow definition
├── shared/                     # Common types
│   ├── types.go               # Request/response types
│   └── utils.go               # Helper functions
└── sample-app/                 # Demo application
    ├── Dockerfile             # Multi-stage Go app
    └── main.go                # Simple HTTP server
```


## Getting Help

- **Temporal Documentation**: https://docs.temporal.io
- **Temporal Community**: https://temporal.io/slack
- **Workshop Issues**: Create an issue in this repository

---

This Part 1 implementation demonstrates the fundamental building blocks of Temporal workflows through a practical CI/CD pipeline example.