# Part 3: Durable Timers - Temporal's Time-Based Workflow Features

This demo showcases Temporal's durable timer capabilities in a real CI/CD pipeline. You'll see how timers survive worker crashes, handle signals, and provide reliable time-based workflow logic.

## What You'll Learn

- **Durable Timers**: Timers that persist across worker restarts and infrastructure failures
- **Selector Patterns**: Waiting for multiple events (timer OR signal) in a single workflow
- **Signal Handling**: Manual intervention to cancel or modify timer behavior
- **Production Safety**: Real-world patterns for deployment validation windows

## Overview

The unified `CICDPipelineWorkflow` now includes all features from Parts 1-3:

```
CICDPipelineWorkflow
├── Docker Build, Test, Push
├── Deploy to Staging (staging/production environments)
├── Human Approval Gate (production environment only)
├── Deploy to Production (production environment only)
└── Durable Timer + Validation (production environment only)
    ├── 30-second countdown timer
    ├── Validation signal can cancel timer
    └── Timer demonstrates durability
```

**Key Feature**: The timer automatically activates for any production deployment - no special configuration needed.

## Prerequisites

- Running Temporal server: `temporal server start-dev`
- Go development environment
- Docker installed and running
- Kubernetes cluster (for deployment phases)

## Demo Instructions

### Step 1: Start the Worker

In Terminal 1:
```bash
go run workers/main.go
```

Expected output:
```
Starting Temporal worker for CI/CD Pipeline
Worker listening on task queue: cicd-task-queue
Registered workflows:
  - CICDPipelineWorkflow (unified workflow with all features)
Registered activities:
  - Docker: Build, Test, Push
  - Kubernetes: Deploy, CheckStatus, Rollback, GetServiceURL
  - Approval: SendRequest, LogDecision, SendNotification
  - Monitoring: ValidateDeployment
```

### Step 2: Start a Production Deployment

In Terminal 2:
```bash
go run cmd/starter/main.go \
  -action=create \
  -image=demo-app \
  -tag=1.0.5 \
  -dockerfile=sample-app/Dockerfile \
  -registry=registry.digitalocean.com/ziggys-containers \
  -env=production
```

Expected output:
```
Started workflow:
  WorkflowID: pipeline-XXXXXXXXX-XXXX
  RunID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  Image: demo-app:1.0.5
  Registry: registry.digitalocean.com/ziggys-containers
  Environment: production

View in Temporal UI: http://localhost:8233/namespaces/default/workflows/pipeline-XXXXXXXXX-XXXX
```

**Important**: Copy the workflow ID for subsequent commands.

### Step 3: Monitor Workflow Progress

Open Temporal UI: http://localhost:8233

Watch the workflow progress through phases:
1. ✅ Docker Build (~1 second)
2. ✅ Docker Test (~3 seconds)  
3. ✅ Docker Push (~1 second)
4. ✅ Staging Deployment (~10 seconds)
5. ⏸️ **Waiting for Approval** (human gate)

### Step 4: Approve Production Deployment

In Terminal 2 (replace `pipeline-XXXXX` with your workflow ID):
```bash
go run cmd/starter/main.go \
  -action=approve \
  -workflow=pipeline-XXXXX \
  -approver="ops-team" \
  -reason="Validated in staging environment"
```

Expected output:
```
✅ Approval signal sent successfully!
  Workflow ID: pipeline-XXXXX
  Approved by: ops-team
  Reason: Validated in staging environment
```

### Step 5: Watch the Timer Phase

After approval, the workflow will:
1. ✅ Deploy to Production (~10 seconds)
2. ⏰ **Start 30-second durable timer**

**In Temporal UI**, you'll see:
- Timer event with ID and 30-second duration
- Real-time countdown in the event history
- Workflow status showing "Running" while timer is active

**In Worker Logs**, you'll see:
```
INFO Phase 5: Starting durable timer demonstration (30 seconds)
INFO This timer demonstrates Temporal's durability - it will survive worker crashes
DEBUG NewTimer TimerID XX Duration 30s
```

### Step 6: Choose Your Demo Path

#### Option A: Validate the Deployment (Cancel Timer)

Send validation signal before timer expires:
```bash
go run cmd/starter/main.go \
  -action=validate \
  -workflow=pipeline-XXXXX \
  -validator="qa-team" \
  -reason="Load testing passed, monitoring looks good"
```

**Result**:
```
✅ Validation signal sent successfully!
🎉 Deployment validated - rollback timer has been canceled!
```

**In Worker Logs**:
```
INFO Validation signal received - timer will be canceled
INFO Timer was canceled by validation signal - demonstrating signal handling
INFO Durable timer demonstration completed
```

#### Option B: Let Timer Expire Naturally

Don't send any validation signal - wait 30 seconds.

**Result**: Timer completes naturally and logs:
```
INFO Timer expired after 30 seconds - demonstrating durable timer completion
INFO Timer completed naturally - no validation signal received
INFO In a real scenario, this could trigger alerts, notifications, or other actions
```

### Step 7: Crash Resilience Demo (Advanced)

**This demonstrates the "durable" aspect of durable timers:**

1. **Start a new workflow** and approve it to reach the timer phase
2. **Kill the worker** while timer is running:
   ```bash
   # In Terminal 1, press Ctrl+C
   # OR from Terminal 2:
   pkill -f "go run workers/main.go"
   ```
3. **Observe in Temporal UI**: Timer continues counting down with no worker
4. **Restart the worker**: `go run workers/main.go`
5. **Timer continues**: Resumes from exact point with no time lost
6. **Send validation**: Signal still works after worker restart

## Key Demonstrations

### 1. Timer Durability
- Timers survive worker crashes and restarts
- No external cron jobs or schedulers needed
- Perfect timing accuracy without clock drift
- State preserved in Temporal server

### 2. Selector Pattern
- Single workflow waits for multiple events
- Timer OR signal - whichever comes first
- Clean cancellation when validation received
- Deterministic behavior every time

### 3. Production Safety Pattern
- Validation windows for critical deployments
- Human oversight with automatic fallback
- Audit trail of who validated what and when
- Real-world deployment safety mechanisms

### 4. Signal Integration
- Real-time manual intervention
- Immediate feedback to operators
- Cancels timer and provides confirmation
- Demonstrates reactive workflow patterns

## Common Demo Scenarios

### Quick Validation (Happy Path)
```bash
# 1. Start workflow
go run cmd/starter/main.go -action=create -image=demo-app -tag=1.0.6 \
  -dockerfile=sample-app/Dockerfile \
  -registry=registry.digitalocean.com/ziggys-containers -env=production

# 2. Approve quickly
go run cmd/starter/main.go -action=approve -workflow=<id> \
  -approver="ops-team" -reason="Quick deployment"

# 3. Validate immediately after production deployment
go run cmd/starter/main.go -action=validate -workflow=<id> \
  -validator="qa-team" -reason="Fast validation"
```

### Forgotten Deployment (Timer Expiration)
```bash
# 1. Start and approve workflow
# 2. Walk away - let timer run for full 30 seconds
# 3. Timer expires naturally demonstrating automatic behavior
```

### Last-Second Validation
```bash
# 1. Start and approve workflow
# 2. Wait ~25 seconds (watch countdown in UI)
# 3. Send validation at last second to show precise timing
```

## Troubleshooting

### Timer Not Appearing
- Ensure environment is set to "production"
- Verify workflow reached production deployment phase
- Check worker logs for "Phase 5: Starting durable timer"
- Look for "NewTimer" debug message in logs

### Signal Not Working
- Verify exact workflow ID (copy from workflow creation output)
- Ensure workflow is still running (not completed)
- Check signal names are exact: "approval", "validation"
- Validate all required parameters are provided

### Activities Stuck
- Check Kubernetes cluster connectivity
- Verify Docker registry access
- Monitor worker heartbeats in Temporal UI
- Consider infrastructure issues if activities timeout

### Worker Crashes During Demo
- **This is a feature!** Shows timer durability
- Restart worker to demonstrate resilience
- Timer continues from exact point it left off
- Use as opportunity to highlight Temporal's reliability

## Understanding Temporal UI

During timer phase, look for:

1. **Event History**: 
   - "TimerStarted" event with timer ID and duration
   - Real-time countdown in event details
   - "TimerCanceled" or "TimerFired" completion events

2. **Workflow Status**:
   - "Running" while timer is active
   - "Completed" after validation or timer expiration

3. **Signal Events**:
   - "SignalReceived" events for approval and validation
   - Signal payload showing approver/validator details

## Architecture Notes

### Why This Pattern Works

1. **No External Dependencies**: Timer state managed entirely by Temporal
2. **Crash Resilient**: Infrastructure failures don't affect timing
3. **Exactly-Once Execution**: Timer fires exactly once when it should
4. **Integrated Workflow State**: Timer and signals share workflow context

### Real-World Applications

- **Deployment Windows**: Automatic rollback if not validated
- **SLA Monitoring**: Alert if process takes too long
- **Approval Timeouts**: Escalate if human doesn't respond
- **Circuit Breakers**: Open/close based on time intervals
- **Batch Processing**: Schedule recurring operations

## Key Takeaways

1. **Durable timers are first-class workflow primitives** in Temporal
2. **Complex time-based logic becomes simple** with workflow selectors
3. **Production safety patterns are easy to implement** without external orchestration
4. **Crash resilience extends to all workflow state** including timers
5. **Real-time observability** shows exactly what's happening when

This demonstrates why Temporal excels for platform engineering: reliable, observable, time-sensitive processes with built-in durability guarantees.

## What's Next

This unified workflow now contains all the CI/CD features you need. You can:
- Extend with additional environments
- Add more validation steps
- Implement different timer durations
- Create custom approval workflows
- Build monitoring and alerting on top

The foundation is production-ready for real CI/CD pipelines.