# Demo 4: Crash Resilience

This demo showcases Temporal's core value proposition: **workflows survive crashes and resume exactly where they left off**. Using the same code as Demo 1, this demonstration focuses on Temporal's durability and deterministic replay capabilities.

## Overview

This demo demonstrates:
- **Workflow Durability**: Workflows survive worker crashes
- **Deterministic Replay**: Workflows resume from exact point of interruption
- **No Duplicate Operations**: Side effects are not repeated
- **State Preservation**: All workflow context is maintained across restarts

## Setup

### Using the Demo4 Tag

```bash
# Switch to the demo4 tag
git checkout demo4-crash-resilience

# Verify you're on the correct code
git describe --tags
# Should output: demo4-crash-resilience
```

This tag points to the Demo 1 codebase, which provides the perfect foundation for crash demonstration due to its:
- **Long-running activities**: Docker build operations take time
- **Side effects**: Actual Docker images are created
- **Clear state**: Easy to observe workflow progress

## Running the Crash Demo

### Prerequisites

1. **Temporal server running**:
   ```bash
   temporal server start-dev
   ```

2. **Clean environment** (remove any existing demo images):
   ```bash
   docker rmi demo-app:crash-test 2>/dev/null || true
   ```

### Step 1: Start the Worker

```bash
go run workers/main.go
```

You should see:
```
Starting Temporal worker for Demo 1 - Basic Pipeline
Worker listening on task queue: cicd-task-queue
```

**Key Point**: Keep this terminal visible to demonstrate the crash.

### Step 2: Start the Workflow

In a **separate terminal**, trigger a workflow:

```bash
go run cmd/starter/main.go -action=create -image=demo-app -tag=crash-test
```

Expected output:
```
Started workflow:
  WorkflowID: pipeline-1640995200-1234
  RunID: a1b2c3d4-e5f6-7890-abcd-ef1234567890
  Image: demo-app:crash-test
  Registry: 
  Environment: staging

View in Temporal UI: http://localhost:8233/namespaces/default/workflows/pipeline-1640995200-1234
```

### Step 3: Monitor Workflow Progress

1. **Open Temporal UI**: http://localhost:8233
2. **Navigate to the workflow** using the WorkflowID from step 2
3. **Watch the workflow progress** through activities

### Step 4: Execute the Crash (Timing is Critical!)

**Wait for the right moment**: 
- Watch the worker logs for "Starting Docker build"
- **During the Docker build activity** (you'll see docker build output)
- **Before the build completes**

**Kill the worker**:

**Option A** - Kill with Ctrl+C:
```bash
# In the worker terminal, press Ctrl+C
^C
```

**Option B** - Kill process from another terminal:
```bash
# Find and kill the worker process
pkill -f "go run workers/main.go"
```

### Step 5: Observe the Crash Effects

**In Temporal UI**:
1. **Refresh the workflow page**
2. **Notice**: Workflow shows "Running" status (not failed!)
3. **Observe**: Current activity shows as "Scheduled" or "Started"
4. **Key point**: Workflow is waiting for a worker, not failing

**In Docker**:
```bash
# Check if partial Docker image exists
docker images | grep demo-app

# Check build cache
docker builder prune --dry-run
```

### Step 6: Restart the Worker (The Magic Moment)

```bash
go run workers/main.go
```

**Immediately observe**:
1. **Worker reconnects** to Temporal
2. **Workflow resumes** from the exact point of interruption
3. **No duplicate Docker build** - may use cache or continue
4. **Workflow completes** successfully

### Step 7: Verify No Duplication

**Check the logs**:
- Look for "Docker build completed" - should only appear once
- Verify image was built successfully
- Confirm no duplicate side effects occurred

**Verify final state**:
```bash
# Image should exist and be properly tagged
docker images | grep demo-app

# Should see single image with crash-test tag
```

## Key Demonstration Points

### 1. Workflow State Preservation

**Before crash**:
- Workflow variables and context maintained
- Activity inputs/outputs preserved
- Timer states saved

**After restart**:
- Exact same workflow state restored
- Variables contain same values
- No state loss occurred

### 2. Deterministic Replay

**Critical concept**: 
- Workflows are deterministic functions
- All non-deterministic operations (like Docker commands) happen in activities
- Activities may retry, but workflow logic replays identically

**Show in UI**:
- Event history shows exact sequence
- No gaps or duplicates in event log
- Deterministic replay ensures consistency

### 3. Activity Idempotency

**Docker build activity**:
- May restart from beginning (safe due to Docker layer caching)
- Or continue from checkpoint (depending on implementation)
- Result is identical regardless

### 4. No Lost Work

**Key value proposition**:
- Long-running CI/CD pipelines can take hours
- Worker crashes, network issues, deployments don't lose progress
- Workflows resume exactly where they left off

## Advanced Crash Scenarios

### Scenario 1: Multiple Crashes

1. Start workflow
2. Crash during Docker build
3. Restart worker
4. Crash again during Docker test
5. Restart worker again
6. Observe: workflow completes successfully

### Scenario 2: Crash During Activity Retry

1. Enable failure simulation: `export SIMULATE_DOCKER_FAILURE=true`
2. Start workflow (will retry due to simulated failures)
3. Crash worker during retry attempts
4. Restart worker
5. Observe: retry count preserved, continues retrying

### Scenario 3: Long-Duration Crash

1. Start workflow
2. Crash worker
3. **Wait several minutes** (demonstrate patience)
4. Restart worker
5. Observe: workflow immediately resumes

## Troubleshooting

### Worker Doesn't Reconnect

**Symptoms**: Restarted worker doesn't pick up workflow

**Solutions**:
```bash
# Check Temporal server is running
temporal operator namespace list

# Verify task queue
temporal workflow list --task-queue cicd-task-queue

# Check worker logs for connection errors
```

### Workflow Shows as Failed

**Symptoms**: Workflow status shows "Failed" instead of "Running"

**Likely cause**: Crash happened after workflow completed

**Solution**: Start a new workflow and time the crash better

### Docker Build Completes Too Quickly

**For slower builds** (better crash timing):
```bash
# Use a larger base image or add dependencies
# Modify sample-app/Dockerfile temporarily
```

**Alternative**: Add artificial delay to Docker build activity for demo purposes

### Can't Find Process to Kill

```bash
# List Go processes
ps aux | grep "go run"

# Or use more specific search
pgrep -f "workers/main.go"
```

## Common Questions & Answers

**Q**: What happens to the Docker build when the worker crashes?
**A**: The Docker daemon continues running independently. When the worker restarts, the activity may use Docker's layer caching or restart the build.

**Q**: How does Temporal know where to resume?
**A**: Temporal stores the complete event history. When a worker reconnects, it replays the workflow deterministically using this history.

**Q**: What if the crash happens between activities?
**A**: Perfect! The workflow resumes at the next activity. Previous activities' results are preserved.

**Q**: Is this specific to Docker builds?
**A**: No. This works for any long-running operation: database migrations, file uploads, API calls, etc.

**Q**: What if multiple workers are running?
**A**: Another worker would pick up the workflow immediately. The crash is invisible to the workflow.

## Integration with Workshop Flow

**Positioning in workshop**:
- Use after Demo 1, 2, or 3
- Provides powerful "wow" moment
- Demonstrates core Temporal value
- No additional code needed

**Timing**: 10-15 minutes
**Impact**: High - often the most memorable demo
**Difficulty**: Low (just timing the crash correctly)

This demonstration effectively showcases why platform engineers choose Temporal for critical infrastructure automation - the guarantee that long-running processes will complete successfully despite inevitable infrastructure failures.