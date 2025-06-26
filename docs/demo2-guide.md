# Demo 2: Human-in-the-Loop CI/CD Pipeline

This demo extends the basic Docker pipeline from Demo 1 to include Kubernetes deployment with human approval gates. It demonstrates how Temporal handles long-running processes with human interaction.

## Key Concepts Demonstrated

1. **Human-in-the-Loop Workflows**: Integration of human decision points in automated processes
2. **Signals**: Using Temporal signals to communicate with running workflows
3. **Kubernetes Integration**: Deploying applications to staging and production environments
4. **Workflow State Preservation**: How workflows maintain state while waiting for approval

## Prerequisites

- Completed Demo 1 setup
- Running Temporal server (`temporal server start-dev`)
- Kubernetes cluster access (Docker Desktop, kind, or cloud provider)
- Kubernetes namespaces created:
  ```bash
  kubectl create namespace staging
  kubectl create namespace production
  ```

## What's New in Demo 2

### New Activities
- **Kubernetes Activities** (`activities/kubernetes.go`):
  - `DeployToKubernetes`: Deploys applications to K8s clusters
  - `CheckDeploymentStatus`: Monitors deployment health
  - `RollbackDeployment`: Rollback capability
  - `GetServiceURL`: Retrieves service endpoints

- **Approval Activities** (`activities/approval.go`):
  - `SendApprovalRequest`: Notifies about pending approvals
  - `LogApprovalDecision`: Records approval/rejection decisions
  - `SendApprovalNotification`: Notifies about deployment decisions

### New Workflow
- **PipelineWithApprovalWorkflow** (`workflows/pipeline.go`):
  - Extends BasicPipelineWorkflow with deployment phases
  - Automatically deploys to staging after successful build/test/push
  - Waits for human approval before production deployment
  - Handles approval/rejection signals

### Enhanced CLI
- **Approval Commands** (`cmd/starter/main.go`):
  - `-action=approve`: Approve a deployment
  - `-action=reject`: Reject a deployment  
  - `-action=status`: Check workflow status

## Running the Demo

### Step 1: Start the Worker

In Terminal 1:
```bash
go run workers/main.go
```

You should see:
```
Starting Temporal worker for Demo 2 - Human-in-the-Loop Pipeline
Worker listening on task queue: cicd-task-queue
Registered workflows:
  - BasicPipelineWorkflow
  - PipelineWithApprovalWorkflow
Registered activities:
  - Docker: Build, Test, Push
  - Kubernetes: Deploy, CheckStatus, Rollback, GetServiceURL
  - Approval: SendRequest, LogDecision, SendNotification
```

### Step 2: Start a Production Deployment

In Terminal 2, trigger a production deployment:
```bash
go run cmd/starter/main.go \
  -action=create \
  -image=demo-app \
  -tag=1.0.0 \
  -dockerfile=sample-app/Dockerfile \
  -registry=registry.digitalocean.com/ziggys-containers \
  -env=production
```

### Step 3: Observe the Workflow Progress

1. Open Temporal UI: http://localhost:8233
2. Find your workflow by ID
3. Watch as it progresses through:
   - Docker build ✓
   - Docker test ✓
   - Docker push ✓
   - Staging deployment ✓
   - **Waiting for approval** ⏸️

### Step 4: Review the Approval Request

The worker logs will show an approval request:
```
==================================================
APPROVAL REQUIRED - Production Deployment
==================================================
Workflow ID: pipeline-XXXXX
Image Tag: demo-app:v2.0.0
Environment: production
Staging URL: http://staging.demo-app.local:8080

To approve:
  go run cmd/starter/main.go -action=approve -workflow=pipeline-XXXXX

To reject:
  go run cmd/starter/main.go -action=reject -workflow=pipeline-XXXXX
==================================================
```

### Step 5: Check Workflow Status

Check the current status:
```bash
go run cmd/starter/main.go \
  -action=status \
  -workflow=pipeline-XXXXX
```

### Step 6: Approve or Reject the Deployment

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

### Step 7: Observe the Result

- **If Approved**: Workflow continues with production deployment
- **If Rejected**: Workflow completes with rejection message

## Key Demo Points

### 1. Long-Running Workflows
- Show how the workflow waits indefinitely for approval
- Stop and restart the worker - workflow continues waiting
- Temporal preserves the entire workflow state

### 2. Signal Handling
- Demonstrate sending signals to running workflows
- Show signal details in Temporal UI
- Explain how signals enable external interaction

### 3. Environment-Based Logic
- Staging deployments happen automatically
- Production requires explicit approval
- Show how workflow adapts based on environment

### 4. Workflow Queries (Advanced)
- While waiting for approval, you can query workflow state
- Useful for building dashboards or status pages

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

# Quick approval with reason
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

## Troubleshooting

### Workflow Not Receiving Signal
- Ensure workflow ID is correct
- Check that workflow is still running (not completed)
- Verify signal name matches ("approval")

### Real Kubernetes Deployments
- Activities now perform actual Kubernetes deployments using kubectl
- Creates real deployments and services in the cluster
- Returns actual service URLs (LoadBalancer or NodePort)
- Requires kubectl access to your cluster

### Activity Timeouts
- Default timeout is 10 minutes per activity
- Adjust `StartToCloseTimeout` for longer operations
- Add heartbeats for very long activities

## What's Next?

In Demo 3, we'll add:
- Deployment windows (business hours only)
- Automatic rollback timers
- Health monitoring and validation
- Advanced selector patterns

## Code Structure

```
activities/
├── docker.go        # From Demo 1
├── kubernetes.go    # NEW: K8s deployment operations
└── approval.go      # NEW: Human approval handling

workflows/
└── pipeline.go      # Extended with PipelineWithApprovalWorkflow

cmd/starter/
└── main.go         # Enhanced with approval commands

workers/
└── main.go         # Registers all activities and workflows
```

## Key Takeaways

1. **Temporal makes human-in-the-loop workflows simple**: Just wait for a signal
2. **State is preserved automatically**: No database or queue management needed
3. **Workflows are durable**: Survive worker crashes, restarts, and delays
4. **Signals enable external interaction**: Perfect for approvals, cancellations, updates
5. **Environment-specific logic is straightforward**: Conditionally execute activities