# Temporal CI/CD Workshop

A comprehensive hands-on workshop demonstrating how Temporal revolutionizes CI/CD pipelines for platform engineering teams through practical, progressive demonstrations.

## Workshop Overview

This workshop consists of multiple demonstrations, each contained in its own branch, showcasing different aspects of Temporal's capabilities:

### Available Branches and Demos

#### Part 1: Basic Docker Pipeline
**Branch**: `part-1-docker-pipeline`
- **What it demonstrates**: Core Temporal concepts with a simple Docker workflow
- **Features**: Build → Test → Push with retry policies and error handling
- **Key learnings**: Workflows as code, activity patterns, observability
- **Quick start**: `git checkout part-1-docker-pipeline`

#### Part 2: Human-in-the-Loop Workflows  
**Branch**: `part-2-human-in-the-loop`
- **What it demonstrates**: Integration of human approval gates in automated workflows
- **Features**: Kubernetes deployment, staging/production environments, approval signals
- **Key learnings**: Long-running workflows, signal handling, state preservation
- **Quick start**: `git checkout part-2-human-in-the-loop`

#### Part 3: Production-Ready Features
**Branch**: `part-3-durable-timers`
- **What it demonstrates**: Advanced production patterns with durable timers
- **Features**: 30-second validation windows, automatic rollbacks, timer vs signal selectors
- **Key learnings**: Durable timers, production safety, selector patterns
- **Quick start**: `git checkout part-3-durable-timers`

#### Part 4: Multi-Language Coordination
**Branch**: `part-4-polyglot-activities`
- **What it demonstrates**: Cross-language worker coordination
- **Features**: Go workflows, Python Kubernetes activities, TypeScript approval activities
- **Key learnings**: Polyglot architecture, specialized workers, task queue routing
- **Quick start**: `git checkout part-4-polyglot-activities`

## Prerequisites

- **Docker Desktop** with Kubernetes enabled
- **Go 1.21+**
- **Temporal CLI** (`brew install temporal`)
- **kubectl** configured for your cluster
- **Python 3.8+** (for Part 4)
- **Node.js 18+** (for Part 4)

## Quick Start (Any Demo)

### 1. Start Temporal Server
```bash
temporal server start-dev
```

### 2. Setup Kubernetes (Required for Parts 2-4)
```bash
./setup/setup-k8s.sh
```

### 3. Choose Your Demo
```bash
# Part 1: Basic Docker Pipeline
git checkout part-1-docker-pipeline

# Part 2: Human Approval Workflows
git checkout part-2-human-in-the-loop

# Part 3: Production Features with Timers
git checkout part-3-durable-timers

# Part 4: Multi-Language Workers
git checkout part-4-polyglot-activities
```

### 4. Start Worker(s)
```bash
# Parts 1-3: Single Go worker
go run workers/main.go

# Part 4: Multiple workers (requires 3 terminals)
# Terminal 1: Go worker
go run workers/main.go

# Terminal 2: Python worker  
cd python-activities && uv run python src/worker.py

# Terminal 3: TypeScript worker
cd typescript-activities && npm start
```

### 5. Run a Pipeline
```bash
# Basic pipeline (Part 1)
go run cmd/starter/main.go -image=demo-app -tag=v1.0.0 -registry=YOUR_REGISTRY

# Production pipeline with approval (Parts 2-4)
go run cmd/starter/main.go -action=create -image=demo-app -tag=v1.0.0 -registry=YOUR_REGISTRY -env=production
```

## Branch Navigation Guide

Each branch has its own README with specific instructions:

| Branch | Purpose | Complexity | Prerequisites |
|--------|---------|------------|---------------|
| `part-1-docker-pipeline` | Basic Docker workflow demonstration | Beginner | Docker, Go |
| `part-2-human-in-the-loop` | Human approval integration | Intermediate | + Kubernetes |
| `part-3-durable-timers` | Production features and timers | Intermediate | + kubectl access |
| `part-4-polyglot-activities` | Multi-language coordination | Advanced | + Python, Node.js |

## Common Commands

### Workflow Operations
```bash
# Create deployment
go run cmd/starter/main.go -action=create -image=<name> -tag=<version> -env=<staging|production>

# Approve production deployment (Parts 2-4)
go run cmd/starter/main.go -action=approve -workflow=<id> -approver="<name>" -reason="<reason>"

# Validate deployment (Parts 3-4)
go run cmd/starter/main.go -action=validate -workflow=<id> -validator="<name>" -reason="<reason>"

# Check status
go run cmd/starter/main.go -action=status -workflow=<id>
```

### Infrastructure
```bash
# Setup Kubernetes namespaces
./setup/setup-k8s.sh

# Check deployments
kubectl get deployments -n staging
kubectl get deployments -n production

# Cleanup resources
./setup/cleanup.sh
```

## Demo Progression

We recommend following the parts in order to build understanding:

1. **Start with Part 1** - Learn core Temporal concepts with a simple workflow
2. **Progress to Part 2** - Add complexity with Kubernetes and human approval
3. **Explore Part 3** - See production patterns with timers and rollbacks  
4. **Finish with Part 4** - Experience multi-language coordination

Each branch builds conceptually on the previous, though they can be run independently.

## Key Concepts Demonstrated

- **Reliability**: Workflows survive crashes and resume exactly where they left off
- **Human Integration**: Seamlessly blend automation with human approval gates
- **Observability**: Complete visibility into long-running processes
- **Durability**: Timers and schedules that survive service restarts
- **Polyglot**: Coordinate activities across multiple programming languages
- **Production Patterns**: Approval gates, rollbacks, validation windows

## Architecture Overview

The workshop demonstrates a unified `CICDPipelineWorkflow` that adapts based on branch:

```
CICDPipelineWorkflow
├── Docker Build, Test, Push (All parts)
├── Deploy to Staging (Parts 2-4)
├── Human Approval Gate (Parts 2-4, production only)
├── Deploy to Production (Parts 2-4, after approval)
└── Validation Timer (Parts 3-4, with rollback)
```

## Getting Help

- **Temporal Documentation**: https://docs.temporal.io
- **Workshop Issues**: Create an issue in this repository
- **Temporal Community**: https://temporal.io/slack

## Clean Up

After exploring the workshop:

```bash
./setup/cleanup.sh
```

---

**Note**: Each branch contains its own detailed README with specific setup instructions, demo scripts, and learning objectives. Switch to any branch and check its README for complete guidance.