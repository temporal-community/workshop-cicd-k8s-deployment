# Temporal CI/CD Workshop - Instructor Guide

This guide provides step-by-step instructions for delivering the Temporal CI/CD workshop.

## Pre-Workshop Setup (15 minutes before)

### 1. Environment Verification

```bash
# Start Temporal development server
temporal server start-dev

# In another terminal, verify Temporal is running
temporal operator namespace list

# Setup Kubernetes namespaces
cd workshop-cicd-k8s-deployment
./setup/setup-k8s.sh

# Verify namespaces
kubectl get ns | grep -E "(staging|production)"
```

### 2. Configure Container Registry

Set your container registry URL as an environment variable:
```bash
export REGISTRY_URL="your-registry-url"  # e.g., registry.digitalocean.com/workshop
```

### 3. Build Base Image

```bash
cd sample-app
docker build -t demo-app:base .
cd ..
```

## Demo Flow

### Demo 1: Basic Docker Pipeline (15 minutes)

**Branch**: `demo1-basic-pipeline`

```bash
git checkout demo1-basic-pipeline
```

**Key Points to Emphasize**:
- Temporal workflows as code
- Activity retry policies
- Workflow visibility in UI

**Demo Steps**:

1. Show the workflow code:
   ```bash
   cat workflows/pipeline.go
   ```

2. Start the worker:
   ```bash
   go run workers/main.go
   ```

3. In another terminal, trigger the workflow:
   ```bash
   go run cmd/starter/main.go -action=create -image=demo-app:v1.0.0
   ```

4. Open Temporal UI: http://localhost:8233

5. Show workflow execution, activity retries, and event history

6. Demonstrate failure simulation:
   ```bash
   # Set environment variable to simulate failures
   export SIMULATE_DOCKER_FAILURE=true
   go run cmd/starter/main.go -action=create -image=demo-app:v1.0.1
   ```

### Demo 2: Human Approval Integration (15 minutes)

**Branch**: `demo2-human-approval`

```bash
git checkout demo2-human-approval
```

**Key Points to Emphasize**:
- Human-in-the-loop workflows
- Signals for external interaction
- Long-running workflow state

**Demo Steps**:

1. Show the enhanced workflow with approval:
   ```bash
   cat workflows/pipeline.go | grep -A 20 "WaitForApproval"
   ```

2. Start the worker:
   ```bash
   go run workers/main.go
   ```

3. Trigger deployment to production:
   ```bash
   go run cmd/starter/main.go -action=create -image=demo-app:v2.0.0 -env=production
   ```

4. Show staging deployment:
   ```bash
   kubectl get pods -n staging
   kubectl get svc -n staging
   ```

5. Demonstrate approval:
   ```bash
   # Get workflow ID from UI or starter output
   go run cmd/starter/main.go -action=approve -workflow=<workflow-id>
   ```

6. Watch production deployment complete

### Demo 3: Production Features (15 minutes)

**Branch**: `demo3-production-features`

```bash
git checkout demo3-production-features
```

**Key Points to Emphasize**:
- Durable timers
- Deployment windows
- Automatic rollbacks
- Selector patterns

**Demo Steps**:

1. Show timer implementation:
   ```bash
   cat workflows/pipeline.go | grep -A 30 "rollbackTimer"
   ```

2. Set deployment window (for demo, make it outside current time):
   ```bash
   export DEPLOYMENT_WINDOW_START=22
   export DEPLOYMENT_WINDOW_END=6
   ```

3. Start worker and trigger deployment:
   ```bash
   go run workers/main.go &
   go run cmd/starter/main.go -action=create -image=demo-app:v3.0.0 -env=production
   ```

4. Show workflow waiting for deployment window in UI

5. Override window for demo:
   ```bash
   export DEPLOYMENT_WINDOW_START=0
   export DEPLOYMENT_WINDOW_END=24
   ```

6. Show rollback timer after deployment

7. Validate deployment to cancel rollback:
   ```bash
   go run cmd/starter/main.go -action=validate -workflow=<workflow-id>
   ```

### Demo 4: Crash Resilience (15 minutes)

**Branch**: Use `demo1-basic-pipeline` (via tag)

```bash
git checkout demo4-crash-resilience
```

**Key Points to Emphasize**:
- Temporal's core value proposition
- Deterministic replay
- No duplicate side effects

**Demo Steps**:

1. Start worker in foreground:
   ```bash
   go run workers/main.go
   ```

2. In another terminal, start a long-running build:
   ```bash
   go run cmd/starter/main.go -action=create -image=demo-app:crash-test
   ```

3. During Docker build activity (watch the logs), kill the worker:
   - Press Ctrl+C in the worker terminal
   - Or use: `pkill -f "go run workers/main.go"`

4. Show workflow still running in UI (with pending activity)

5. Restart the worker:
   ```bash
   go run workers/main.go
   ```

6. Watch workflow continue from exact point of interruption

7. Verify no duplicate Docker operations occurred

### Demo 5: Polyglot Finale (15 minutes)

**Branch**: `demo5-polyglot-finale`

```bash
git checkout demo5-polyglot-finale
```

**Key Points to Emphasize**:
- Language-specific strengths
- Unified orchestration
- Cross-language coordination

**Demo Steps**:

1. Show the three workers:
   ```bash
   ls workers/
   ```

2. Start all three workers (in separate terminals):
   ```bash
   # Terminal 1 - Go worker
   go run workers/main.go

   # Terminal 2 - Python worker
   cd workers && python main.py

   # Terminal 3 - Node.js worker
   cd workers && npm install && node main.js
   ```

3. Trigger the polyglot workflow:
   ```bash
   go run cmd/starter/main.go -action=create -image=demo-app:polyglot -env=production
   ```

4. In Temporal UI, show activities executing on different workers

5. Point out language-specific implementations:
   - Go: Infrastructure operations
   - Python: Testing and security scanning
   - Node.js: Notifications and monitoring

## Common Issues and Solutions

### Issue: Container registry authentication
**Solution**: Ensure Docker is logged into the registry:
```bash
docker login $REGISTRY_URL
```

### Issue: Kubernetes permissions
**Solution**: Verify kubectl context:
```bash
kubectl config current-context
kubectl auth can-i create deployments -n staging
```

### Issue: Worker not picking up activities
**Solution**: Check worker registration and task queue names match

### Issue: Temporal server not accessible
**Solution**: Ensure temporal server is running:
```bash
temporal server start-dev
```

## Wrap-up Talking Points

1. **Real-world applications**:
   - Multi-stage deployments
   - Infrastructure provisioning
   - Data pipeline orchestration
   - Batch job processing

2. **Temporal advantages**:
   - Built-in durability
   - Human-in-the-loop patterns
   - Visibility and debugging
   - Language agnostic

3. **Next steps**:
   - Temporal Cloud
   - Production deployment patterns
   - Advanced workflow patterns
   - Custom UI integration

## Cleanup

After the workshop:
```bash
./setup/cleanup.sh
temporal operator namespace delete temporal-demo
```