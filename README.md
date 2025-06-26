# Temporal CI/CD Workshop - Kubernetes Deployment Pipeline

A comprehensive workshop demonstrating Temporal's capabilities for building resilient CI/CD pipelines with Docker and Kubernetes. Features progressive complexity across multiple demo branches showcasing real-world deployment patterns.

## Quick Start

```bash
# Start Temporal server
temporal server start-dev

# Start worker
go run workers/main.go

# Trigger production pipeline
go run cmd/starter/main.go -action=create -image=demo-app -tag=1.0.0 \
  -dockerfile=sample-app/Dockerfile \
  -registry=YOUR_REGISTRY \
  -env=production
```

## Architecture Overview

The workshop implements a unified `CICDPipelineWorkflow` that demonstrates:

- **Docker Operations**: Build, test, and push images with retry policies
- **Kubernetes Deployment**: Automated staging and production deployments  
- **Human Approval Gates**: Production deployments require manual approval
- **Durable Timers**: 30-second validation windows with automatic rollback
- **Polyglot Activities**: Multi-language coordination (Go, TypeScript, Python)
- **Crash Resilience**: Workflows survive worker failures and infrastructure issues

### Workflow Phases

```
CICDPipelineWorkflow
├── Docker Build, Test, Push (Go)
├── Deploy to Staging (Automatic)
├── Human Approval Gate (Production only)
├── Deploy to Production (After approval)
└── Validation Timer (30s rollback window)
```

## Demo Progression

### Demo 1: Basic Docker Pipeline
**Branch**: `demo1-basic-pipeline`

- Sequential Docker operations (build → test → push)
- Retry policies with exponential backoff
- Activity timeouts and error simulation
- Comprehensive logging and observability

**Key Commands**:
```bash
# Basic pipeline
go run cmd/starter/main.go -image=demo-app -tag=v1.0.0

# Simulate failures
export SIMULATE_DOCKER_FAILURE=true
go run cmd/starter/main.go -image=demo-app -tag=v1.0.1
```

### Demo 2: Human-in-the-Loop Pipeline
**Branch**: `demo2-human-approval`

- Kubernetes deployment activities
- Human approval workflows with signals
- Environment-based deployment logic
- Long-running workflow state preservation

**Key Commands**:
```bash
# Production deployment requiring approval
go run cmd/starter/main.go -action=create -image=demo-app -tag=1.0.0 -env=production

# Approve deployment
go run cmd/starter/main.go -action=approve -workflow=<id> \
  -approver="ops-team" -reason="Validated in staging"

# Check status
go run cmd/starter/main.go -action=status -workflow=<id>
```

### Demo 3: Production-Ready Features
**Branch**: `demo3-production-features`

- Durable timers for validation windows
- Automatic rollback capabilities
- Selector patterns (timer OR signal)
- Production safety mechanisms

**Key Commands**:
```bash
# Validate deployment (cancels timer)
go run cmd/starter/main.go -action=validate -workflow=<id> \
  -validator="qa-team" -reason="Load testing passed"

# Let timer expire for automatic behavior
# (Wait 30 seconds - no validation command)
```

### Demo 4: Crash Resilience
**Branch**: `demo4-crash-resilience`

- Worker crash simulation during various phases
- State durability demonstration
- Timer persistence across restarts
- Recovery without data loss

**Key Commands**:
```bash
# Kill worker during execution
pkill -f "go run workers/main.go"

# Restart worker - workflow continues
go run workers/main.go
```

### Part 4: Polyglot Activities
**Branch**: `part-4-polyglot-activities`

- Multi-language worker coordination
- TypeScript approval activities
- Python Kubernetes activities
- Go workflow orchestration

**Setup Commands**:
```bash
# Terminal 1: Go worker (workflows + Docker)
go run workers/main.go

# Terminal 2: TypeScript worker (approvals)
cd typescript-activities && npm start

# Terminal 3: Python worker (Kubernetes)
cd python-activities && uv run python src/worker.py
```

## Core Commands

### Workflow Operations
```bash
# Create deployment
go run cmd/starter/main.go -action=create -image=<name> -tag=<version> -env=<staging|production>

# Approve production deployment
go run cmd/starter/main.go -action=approve -workflow=<id> -approver="<name>" -reason="<reason>"

# Reject deployment
go run cmd/starter/main.go -action=reject -workflow=<id> -approver="<name>" -reason="<reason>"

# Validate deployment (cancel timer)
go run cmd/starter/main.go -action=validate -workflow=<id> -validator="<name>" -reason="<reason>"

# Check workflow status
go run cmd/starter/main.go -action=status -workflow=<id>
```

### Infrastructure Setup
```bash
# Setup Kubernetes namespaces
./setup/setup-k8s.sh

# Cleanup resources
./setup/cleanup.sh

# Check deployments
kubectl get deployments -n staging
kubectl get deployments -n production
```

### Kubernetes Requirements
- `staging` and `production` namespaces
- kubectl configured with cluster access
- LoadBalancer or NodePort service support

## Learning Outcomes

After completing this workshop, you'll understand:

1. **Workflow Durability**: How Temporal maintains state across failures
2. **Activity Patterns**: Retry policies, timeouts, and heartbeats
3. **Signal Handling**: External communication with running workflows
4. **Timer Management**: Time-based workflow logic and validation windows
5. **Polyglot Architecture**: Multi-language activity coordination
6. **Production Patterns**: Approval gates, rollbacks, and operational safety

This workshop provides a production-ready foundation for implementing reliable CI/CD pipelines with Temporal's workflow orchestration capabilities.