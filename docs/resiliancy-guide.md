# Part 4: Crash Resilience Demo

This demo demonstrates how Temporal workflows survive worker crashes and resume exactly where they left off.

## Prerequisites

- Running Temporal server: `temporal server start-dev`
- Go development environment
- Docker installed and running

## Demo Scenarios

### Scenario 1: Crash During Docker Build

#### Step 1: Start Worker and Workflow

Terminal 1 - Start worker:
```bash
go run workers/main.go
```

Terminal 2 - Start workflow:
```bash
go run cmd/starter/main.go \
  -action=create \
  -image=demo-app \
  -tag=crash-test-1 \
  -dockerfile=sample-app/Dockerfile \
  -registry=registry.digitalocean.com/ziggys-containers \
  -env=staging
```

Copy the workflow ID from the output.

#### Step 2: Kill Worker During Build

Watch worker logs for:
```
INFO    Starting Docker build
INFO    Executing docker build
```

Kill the worker immediately:
```bash
# Option 1: In Terminal 1, press Ctrl+C
# Option 2: From Terminal 3:
pkill -f "go run workers/main.go"
```

#### Step 3: Check Temporal UI

Open http://localhost:8233 and find your workflow:
- Status shows "Running" (not failed)
- Pending Activities shows "BuildDockerImage"
- History shows activity started but not completed

#### Step 4: Restart Worker

Terminal 1:
```bash
go run workers/main.go
```

Observe:
- Worker picks up pending activity
- Docker build completes
- Workflow continues to next activity

### Scenario 2: Crash Between Activities

#### Step 1: Start workflow (same as above)

#### Step 2: Kill After Build Completes

Watch for:
```
INFO    Docker build completed
```

Kill worker immediately after this message appears.

#### Step 3: Check Temporal UI

- Build activity shows "Completed"
- No test activity started yet
- Workflow status is "Running"

#### Step 4: Restart Worker

Worker continues from test activity without re-executing build.

### Scenario 3: Production Workflow with Multiple Activities

Start a production workflow for more crash opportunities:

```bash
go run cmd/starter/main.go \
  -action=create \
  -image=demo-app \
  -tag=multi-crash \
  -dockerfile=sample-app/Dockerfile \
  -registry=registry.digitalocean.com/ziggys-containers \
  -env=production
```

Crash at any of these points:
1. During Docker build
2. During staging deployment  
3. During approval wait
4. After approving, during production deployment

Each time:
1. Kill worker
2. Check Temporal UI (workflow still running)
3. Restart worker
4. Workflow continues from exact point

### Scenario 4: Crash During Timer Phase

If using production environment:

1. Start workflow and approve it to reach timer phase
2. Kill worker while 30-second timer is running
3. Check Temporal UI - timer continues counting down
4. Restart worker
5. Timer resumes from exact second
6. Can still send validation signal

### Scenario 5: Extended Downtime

1. Start any workflow
2. Kill worker during any activity
3. Wait 5+ minutes
4. Check Temporal UI - workflow still "Running"
5. Restart worker - immediate recovery

## Commands Reference

### Start Worker
```bash
go run workers/main.go
```

### Kill Worker Options
```bash
# Graceful stop (Ctrl+C in worker terminal)

# Force kill from another terminal
pkill -f "go run workers/main.go"
```

### Check Workflow Status
```bash
go run cmd/starter/main.go -action=status -workflow=<workflow-id>
```

### Approve Production Deployment
```bash
go run cmd/starter/main.go \
  -action=approve \
  -workflow=<workflow-id> \
  -approver="ops-team" \
  -reason="Validated in staging"
```

### Validate Deployment (Cancel Timer)
```bash
go run cmd/starter/main.go \
  -action=validate \
  -workflow=<workflow-id> \
  -validator="qa-team" \
  -reason="Deployment validated"
```

## What to Observe

### In Temporal UI During Crash
- Workflow status remains "Running"
- Pending Activities tab shows waiting activities
- Event history shows incomplete activity

### After Worker Restart
- Pending activities immediately resume
- No duplicate executions
- Workflow continues normally

### Key Demonstrations
1. Workflows never fail due to worker crashes
2. No lost progress or duplicate work
3. Automatic recovery without manual intervention
4. All workflow state (activities, timers, signals) is durable