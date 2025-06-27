# Temporal CI/CD Workshop - Human-in-the-Loop Pipeline

A hands-on demonstration of how Temporal enables sophisticated CI/CD pipelines with human approval gates for platform engineering teams.

## What This Demo Shows

This implementation demonstrates a **human-in-the-loop CI/CD pipeline** that seamlessly integrates automation with human decision points:

1. **Phase 1**: Docker build → test → push (fully automated)
2. **Phase 2**: Deploy to staging environment (automatic)
3. **Phase 3**: Human approval gate for production deployment
4. **Phase 4**: Deploy to production (after approval)

## Key Concepts Demonstrated

- **Human-in-the-Loop Workflows**: Integration of human decision points in automated processes
- **Signals**: Using Temporal signals to communicate with running workflows
- **Kubernetes Integration**: Real deployments to staging and production environments
- **Workflow State Preservation**: How workflows maintain state while waiting for approval
- **Long-Running Processes**: Workflows that can wait indefinitely for human input
- **Environment-Based Logic**: Different behavior for staging vs production deployments

## Prerequisites

- **Docker Desktop** with Kubernetes enabled
- **Go 1.21+**
- **Temporal CLI** (`brew install temporal`)
- **kubectl** configured for your cluster
- **Container registry access** (optional for local testing)

## Quick Start

### 1. Start Temporal Server
```bash
temporal server start-dev
```

### 2. Setup Kubernetes Namespaces
```bash
./setup/setup-k8s.sh
```

### 3. Start the Worker
```bash
go run workers/main.go
```

You should see:
```
Starting Temporal worker for CI/CD Pipeline
Worker listening on task queue: cicd-task-queue
Registered workflows:
  - CICDPipelineWorkflow (human-in-the-loop workflow)
Registered activities:
  - Docker: Build, Test, Push
  - Kubernetes: Deploy, CheckStatus, GetServiceURL
  - Approval: SendRequest, LogDecision, SendNotification
```

### 4. Trigger a Production Deployment
```bash
go run cmd/starter/main.go \
  -action=create \
  -image=demo-app \
  -tag=v1.0.0 \
  -env=production
```

### 5. Watch the Workflow Progress

1. **Open Temporal UI**: http://localhost:8233
2. **Find your workflow** by ID
3. **Watch the progress**:
   - Docker build ✓
   - Docker test ✓
   - Docker push ✓
   - Staging deployment ✓
   - **Waiting for approval** ⏸️

### 6. Review the Approval Request

The worker logs will show:
```
==================================================
APPROVAL REQUIRED - Production Deployment
==================================================
Workflow ID: pipeline-XXXXX
Image Tag: demo-app:v1.0.0
Environment: production
Staging URL: http://staging.demo-app.local:8080

To approve:
  go run cmd/starter/main.go -action=approve -workflow=pipeline-XXXXX

To reject:
  go run cmd/starter/main.go -action=reject -workflow=pipeline-XXXXX
==================================================
```

### 7. Approve or Reject the Deployment

#### To Approve:
```bash
go run cmd/starter/main.go \
  -action=approve \
  -workflow=pipeline-XXXXX \
  -approver="john.doe" \
  -reason="Tested successfully in staging"
```

#### To Reject:
```bash
go run cmd/starter/main.go \
  -action=reject \
  -workflow=pipeline-XXXXX \
  -approver="jane.smith" \
  -reason="Found performance issues in staging"
```

## Available Commands

### Pipeline Operations
```bash
# Start staging deployment (no approval needed)
go run cmd/starter/main.go -action=create -image=demo-app -tag=v1.0.0 -env=staging

# Start production deployment (requires approval)
go run cmd/starter/main.go -action=create -image=demo-app -tag=v1.0.0 -env=production

# Check workflow status
go run cmd/starter/main.go -action=status -workflow=<workflow-id>

# Approve deployment
go run cmd/starter/main.go -action=approve -workflow=<workflow-id> -approver="user" -reason="ready"

# Reject deployment
go run cmd/starter/main.go -action=reject -workflow=<workflow-id> -approver="user" -reason="not ready"
```

### Infrastructure Commands
```bash
# Check Kubernetes deployments
kubectl get deployments -n staging
kubectl get deployments -n production

# View service endpoints
kubectl get services -n staging
kubectl get services -n production

# Clean up resources
./setup/cleanup.sh
```

## Project Structure

```
├── README.md                    # This file
├── activities/                  # Business logic activities
│   ├── approval.go             # Human approval workflows
│   ├── docker.go               # Docker build/test/push operations
│   └── kubernetes.go           # K8s deployment operations
├── cmd/starter/main.go         # CLI interface for triggering workflows
├── workers/main.go             # Temporal worker registration
├── workflows/pipeline.go       # Main CI/CD pipeline workflow
├── shared/                     # Common types and utilities
│   ├── types.go               # Request/response types
│   └── utils.go               # Helper functions
├── sample-app/                 # Demo application to deploy
│   ├── Dockerfile             # Multi-stage Go app container
│   ├── main.go                # Simple HTTP server
│   ├── integration_test.go    # Basic health check tests
│   └── k8s/                   # Kubernetes manifests
└── setup/                     # Environment setup scripts
    ├── setup-k8s.sh          # Creates K8s namespaces
    └── cleanup.sh             # Cleanup script
```

## Common Scenarios

### Scenario 1: Quick Staging Test
```bash
# Deploy to staging only (no approval needed)
go run cmd/starter/main.go \
  -action=create \
  -image=demo-app \
  -tag=test-feature \
  -env=staging
```

### Scenario 2: Emergency Production Fix
```bash
# Start deployment
go run cmd/starter/main.go \
  -action=create \
  -image=demo-app \
  -tag=hotfix-v2.0.1 \
  -env=production

# Quick approval
go run cmd/starter/main.go \
  -action=approve \
  -workflow=<workflow-id> \
  -approver="ops-team" \
  -reason="Critical security fix"
```

### Scenario 3: Rejection After Review
```bash
# Reject with detailed reason
go run cmd/starter/main.go \
  -action=reject \
  -workflow=<workflow-id> \
  -approver="qa-team" \
  -reason="Failed load testing requirements"
```

## Demo Highlights

### 1. Long-Running Workflows
- Workflows wait indefinitely for approval without consuming resources
- Stop and restart the worker - workflow continues waiting exactly where it left off
- Temporal preserves the entire workflow state automatically

### 2. Signal-Based Human Integration
- External systems can send approval/rejection signals to running workflows
- Signals are visible in the Temporal UI with full details
- No polling or database coordination required

### 3. Environment-Specific Behavior
- Staging deployments happen automatically after successful build/test
- Production deployments require explicit human approval
- Same workflow code adapts behavior based on environment parameter

### 4. Reliable State Management
- No lost approvals or missed signals
- Complete audit trail of all decisions
- Workflow history shows exact sequence of events

## Architecture Patterns

### Temporal Patterns Used
- **Signal-based coordination** for human-in-the-loop workflows
- **Activity composition** with retry policies and error handling
- **Conditional workflow logic** based on environment parameters
- **Long-running workflow state** preservation

### Infrastructure Integration
- **Real Kubernetes deployments** using kubectl
- **Multi-environment isolation** with namespace separation
- **Service discovery** and health checking
- **Comprehensive logging** for observability

## Troubleshooting

### Workflow Not Receiving Signals
- Ensure workflow ID is correct
- Check that workflow is still running (not completed)
- Verify signal name matches ("approval")

### Kubernetes Deployment Issues
- Verify kubectl access to your cluster
- Check that staging/production namespaces exist
- Ensure sufficient cluster resources for deployments

### Activity Timeouts
- Default timeout is 10 minutes per activity
- Activities include heartbeats for long operations
- Check Temporal UI for activity details and errors

## Key Takeaways

1. **Temporal makes human-in-the-loop workflows simple**: Just wait for a signal
2. **State is preserved automatically**: No database or queue management needed  
3. **Workflows are durable**: Survive worker crashes, restarts, and delays
4. **Signals enable external interaction**: Perfect for approvals, cancellations, updates
5. **Environment-specific logic is straightforward**: Conditionally execute activities

## Getting Help

- **Temporal Documentation**: https://docs.temporal.io
- **Temporal Community**: https://temporal.io/slack
- **Repository Issues**: Create an issue in this repository

---

This implementation demonstrates how Temporal transforms complex CI/CD automation challenges into simple, reliable workflows that seamlessly integrate human oversight with automated processes.