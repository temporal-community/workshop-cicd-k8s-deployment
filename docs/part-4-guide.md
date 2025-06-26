# Part 4: Polyglot Activities - Multi-Language Worker Coordination

This guide demonstrates how to run approval activities in TypeScript and Kubernetes activities in Python, while keeping the main workflow in Go. This showcases Temporal's powerful polyglot capabilities where different activities can be implemented in different programming languages.

## Overview

In this demonstration:
- **Go Worker**: Runs workflows and Docker activities only
- **TypeScript Worker**: Handles approval activities (SendApprovalRequest, LogApprovalDecision, SendApprovalNotification)
- **Python Worker**: Handles Kubernetes activities (DeployToKubernetes, CheckDeploymentStatus, RollbackDeployment, GetServiceURL)

All workers connect to the same Temporal cluster and task queue, allowing seamless inter-language communication.

## Prerequisites

- Go 1.21+ (for the main workflow and Docker activities)
- Node.js 18+ (for TypeScript approval activities)
- Python 3.8+ (for Python Kubernetes activities)
- [uv](https://github.com/astral-sh/uv) package manager for Python
- Temporal server running
- Docker and kubectl configured


## Step 2: Set Up TypeScript Approval Activities

Navigate to the TypeScript activities directory and install dependencies:

```bash
cd typescript-activities
npm install
```

Build the TypeScript code:

```bash
npm run build
```

Run the TypeScript Worker:

```bash
npm start
```

## Step 3: Set Up Python Kubernetes Activities

Navigate to the Python activities directory and install dependencies using uv:

```bash
cd python-activities
uv sync
```

## Step 4: Start the Temporal Server

Start the Temporal development server:

```bash
temporal server start-dev
```

This will start the Temporal server on `localhost:7233` with the Web UI available at `http://localhost:8233`.

## Step 5: Start All Workers

You need to start all three workers in separate terminals. The order doesn't matter as they all connect to the same Temporal cluster.

### Terminal 1: Start Go Worker (Workflows + Docker Activities)

```bash
# Option A: Using the new polyglot worker file
go run workers/polyglot-main.go

# Option B: Using environment variable with existing worker
POLYGLOT_MODE=true go run workers/main.go
```

You should see:
```
Starting Temporal worker for CI/CD Pipeline (Polyglot Mode)
Worker listening on task queue: cicd-task-queue
Registered workflows:
  - CICDPipelineWorkflow (unified workflow with all features)
Registered activities:
  - Docker: Build, Test, Push
  - Kubernetes: Handled by Python worker
  - Approval: Handled by TypeScript worker
```

### Terminal 2: Start TypeScript Worker (Approval Activities)

```bash
cd typescript-activities
npm run dev
```

You should see:
```
Starting TypeScript worker for approval activities
Connecting to Temporal server: localhost:7233
Worker listening on task queue: cicd-task-queue
Registered activities:
  - TypeScript Approval: SendApprovalRequest, LogApprovalDecision, SendApprovalNotification
Worker started successfully!
```

### Terminal 3: Start Python Worker (Kubernetes Activities)

```bash
cd python-activities
uv run python src/worker.py
```

You should see:
```
Starting Python worker for Kubernetes activities
Connecting to Temporal server: localhost:7233
Worker listening on task queue: cicd-task-queue
Registered activities:
  - Python Kubernetes: DeployToKubernetes, CheckDeploymentStatus, RollbackDeployment, GetServiceURL
Worker started successfully!
```

## Step 6: Set Up Kubernetes Environment

Make sure your Kubernetes namespaces are set up:

```bash
./setup/setup-k8s.sh
```

## Step 7: Run a Polyglot Pipeline

Now trigger a pipeline that will use activities from all three languages:

```bash
go run cmd/starter/main.go \
  -action=create \
  -image=demo-app \
  -tag=polyglot-v1.0.0 \
  -dockerfile=sample-app/Dockerfile \
  -registry=your-registry.com/demo \
  -env=production
```

## What Happens During Execution

1. **Go Workflow**: The CICDPipelineWorkflow starts and orchestrates the entire pipeline
2. **Go Docker Activities**: BuildDockerImage, TestDockerContainer, and PushToRegistry execute in Go
3. **Python Kubernetes Activity**: DeployToKubernetes runs in Python for staging deployment
4. **TypeScript Approval Activity**: SendApprovalRequest runs in TypeScript, logging the approval request
5. **Human Approval**: You approve/reject via the CLI
6. **TypeScript Approval Activity**: LogApprovalDecision runs in TypeScript to log the decision
7. **Python Kubernetes Activity**: DeployToKubernetes runs in Python for production deployment
8. **Durable Timer**: 30-second validation timer (handled by Go workflow)
9. **Validation or Rollback**: Python RollbackDeployment may execute if validation timeout occurs

## Monitoring the Polyglot Execution

### Temporal Web UI
Visit `http://localhost:8233` to see:
- Workflow execution with activities from different languages
- Activity task distribution across workers
- Execution timeline showing polyglot coordination

### Worker Logs
Monitor all three terminal windows to see:
- **Go**: Workflow decisions and Docker activity execution
- **TypeScript**: Approval activity execution with structured logging
- **Python**: Kubernetes operations with kubectl command outputs

### Approval Workflow
When the approval request is logged by TypeScript, approve the deployment:

```bash
go run cmd/starter/main.go \
  -action=approve \
  -workflow=<workflow-id> \
  -approver="polyglot-demo" \
  -reason="Testing multi-language coordination"
```

## Verification Commands

Check that activities are running in different languages:

```bash
# Check deployments created by Python worker
kubectl get deployments -n staging
kubectl get deployments -n production

# Check workflow status
go run cmd/starter/main.go -action=status -workflow=<workflow-id>
```

## Troubleshooting

### Worker Connection Issues
- Ensure all workers can connect to `localhost:7233`
- Check that Temporal server is running: `temporal server start-dev`
- Verify no firewall blocking connections

### Missing Activities
- Ensure all three workers are running and healthy
- Check worker logs for registration messages
- Verify task queue name is consistent: `cicd-task-queue`

### Language-Specific Issues

**TypeScript:**
```bash
# Check TypeScript compilation
npm run build

# Install missing dependencies
npm install
```

**Python:**
```bash
# Check Python dependencies
uv sync

# Verify kubectl is available
kubectl version --client
```

**Go:**
```bash
# Check Go modules
go mod tidy

# Verify Docker is available
docker version
```

## Key Benefits Demonstrated

1. **Language Choice**: Use the best language for each type of activity
2. **Team Expertise**: Different teams can contribute in their preferred languages
3. **Ecosystem Integration**: Leverage language-specific libraries and tools
4. **Gradual Migration**: Migrate activities to different languages incrementally
5. **Operational Consistency**: All activities share the same operational model

## Next Steps

- Experiment with adding activities in other languages (Java, C#, etc.)
- Try migrating some Docker activities to other languages
- Implement language-specific error handling and retry policies
- Add language-specific monitoring and observability

This polyglot setup demonstrates Temporal's unique ability to coordinate activities across multiple programming languages while maintaining strong consistency, durability, and observability across the entire workflow.