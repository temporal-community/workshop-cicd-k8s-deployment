# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Core Workflow Operations
- **Start Temporal server**: `temporal server start-dev`
- **Start worker**: `go run workers/main.go`
- **Trigger pipeline**: `go run cmd/starter/main.go -action=create -image=demo-app -tag=1.0.0 -dockerfile=sample-app/Dockerfile -registry=registry.digitalocean.com/ziggys-containers -env=production`

### Approval Workflow Commands
- **Check workflow status**: `go run cmd/starter/main.go -action=status -workflow=<workflow-id>`
- **Approve deployment**: `go run cmd/starter/main.go -action=approve -workflow=<workflow-id> -approver="<name>" -reason="<reason>"`
- **Reject deployment**: `go run cmd/starter/main.go -action=reject -workflow=<workflow-id> -approver="<name>" -reason="<reason>"`
- **Validate deployment**: `go run cmd/starter/main.go -action=validate -workflow=<workflow-id> -validator="<name>" -reason="<reason>"`

### Kubernetes Setup
- **Setup namespaces**: `./setup/setup-k8s.sh`
- **Cleanup resources**: `./setup/cleanup.sh`
- **Check deployments**: `kubectl get deployments -n staging` or `kubectl get deployments -n production`

### Testing and Build
- **Run sample app tests**: `cd sample-app && go test`
- **Build sample app**: `cd sample-app && go build`
- **Run integration tests**: `cd sample-app && BASE_URL=http://localhost:<port> go test`

## Architecture Overview

This is a Temporal CI/CD workshop demonstrating progressive complexity across multiple demo branches. The codebase uses Go with Temporal SDK to implement durable workflows for Docker-based deployment pipelines.

### Key Components

**Workflows** (`workflows/pipeline.go`):
- `CICDPipelineWorkflow`: Unified workflow implementing full CI/CD pipeline with environment-based branching, approval gates, and durable timers for rollback validation

**Activities** (`activities/`):
- **Docker activities** (`docker.go`): Build, test, and push Docker images with multi-architecture support
- **Kubernetes activities** (`kubernetes.go`): Deploy to staging/production namespaces, health checks, rollbacks, service URL retrieval
- **Approval activities** (`approval.go`): Human-in-the-loop approval notifications and logging

**Worker** (`workers/main.go`):
- Registers all workflows and activities on task queue "cicd-task-queue"
- Handles connection to Temporal server (default: localhost:7233)

**CLI** (`cmd/starter/main.go`):
- Multi-action CLI supporting workflow creation, approval/rejection, validation, and status queries
- Handles workflow signals for human approval integration and deployment validation
- Actions: `create`, `approve`, `reject`, `validate`, `status`

### Branch Structure
- `main`: Foundation and setup
- `demo1-basic-pipeline`: Basic Docker workflow
- `demo2-human-approval`: Human-in-the-loop workflows
- `demo3-production-features`: Production-ready features
- `demo4-crash-resilience`: Demonstrates durability
- `part-4-polyglot-activities`: Multi-language coordination (current branch)

### Data Flow
1. CLI triggers `CICDPipelineWorkflow` with `PipelineRequest`
2. Workflow executes Docker activities sequentially (build → test → push)
3. Automatic deployment to staging environment via `DeployToKubernetes` activity
4. For production environments, workflow waits for approval signal
5. Human approver sends signal via CLI, workflow continues or terminates based on decision
6. Production deployment with durable timer for validation (30 seconds)
7. Validation signal cancels timer, or timer expires triggering automatic rollback

### Key Design Patterns
- **Retry policies**: All activities have exponential backoff (max 5 attempts, 10min timeout)
- **Signal handling**: Workflows wait indefinitely for approval signals using `GetSignalChannel`
- **Durable timers**: Production deployments use 30-second validation timers that survive worker crashes
- **Selector pattern**: Workflows use selectors to wait for either validation signals OR timer expiration
- **Environment-based branching**: Staging deployments are automatic, production requires approval
- **Heartbeat recording**: Long-running activities send heartbeats for visibility
- **Idempotent operations**: Activities can be safely retried without side effects

### Configuration
- Temporal server connection configurable via `TEMPORAL_HOST` environment variable
- Docker failure simulation via `SIMULATE_DOCKER_FAILURE=true`
- Registry push failure simulation via `SIMULATE_PUSH_FAILURE=true`
- Kubernetes deployments target `staging` and `production` namespaces

### Dependencies
- Temporal SDK v1.26.1
- Docker with buildx for multi-architecture builds
- kubectl for Kubernetes operations
- Go 1.21+ for all components