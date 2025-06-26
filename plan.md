# Temporal CI/CD Workshop - Complete Implementation Plan

## Workshop Overview

### Goals
- **Primary**: Demonstrate Temporal's value for platform engineering through realistic CI/CD scenarios
- **Secondary**: Show how Temporal handles crash resilience, human approval, timers, and polyglot workflows
- **Outcome**: Attendees understand how to apply Temporal to their own infrastructure automation

### Target Audience
Platform engineers familiar with Docker, Kubernetes, and CI/CD concepts but new to Temporal

### Workshop Format
- **Duration**: 90 minutes
- **Format**: Live coding demonstrations with progressive complexity
- **Interaction**: Attendees observe, Q&A throughout, take away working code repository

### Core Value Propositions to Demonstrate
1. **Reliability**: Workflows survive crashes and resume exactly where they left off
2. **Human Integration**: Seamlessly blend automation with human approval gates
3. **Observability**: Complete visibility into long-running processes
4. **Polyglot**: Coordinate activities across multiple programming languages
5. **Durability**: Timers and schedules that survive service restarts

### Tech Stack
- The container registry will be DigitalOcean's container registry
- The Kubernetes cluster will be DigitalOcean's kuberenetes offering

---

## Project Structure & Git Branching Strategy

### Final Project Structure when completed
```
temporal-cicd-workshop/
├── README.md                           # Quick start with branch navigation
├── setup/
│   ├── setup-k8s.sh                   # Kubernetes namespace setup
│   └── cleanup.sh                     # Cleanup script
├── sample-app/                        # Application to deploy
│   ├── Dockerfile
│   ├── main.go                        # Simple HTTP server
│   ├── main_test.go                   # Unit tests
│   ├── k8s/
│   │   ├── staging-deployment.yaml
│   │   ├── production-deployment.yaml
│   │   └── service.yaml
│   └── tests/
│       └── integration_test.py
├── shared/                            # Common types and utilities
│   ├── types.go                       # Workflow/activity request/response types
│   └── utils.go                       # Helper functions
├── workflows/                         # Evolves across branches
│   └── pipeline.go                    # Main workflow (changes per part)
├── activities/                        # Grows with each branch
│   ├── docker.go                      # Present in all branches
│   ├── kubernetes.go                  # Added in part-2+
│   ├── approval.go                    # Added in part-2+
│   ├── monitoring.go                  # Added in part-3+
│   ├── python/                        # Added in part-5
│   │   ├── test_runner.py
│   │   └── security_scan.py
│   └── nodejs/                        # Added in part-5
│       ├── notifications.js
│       └── monitoring.js
├── python-code/                       # Python Activities and Worker implementation for polyglot (part-5 only)
│   ├── activities.py                  # Python Activities (part-5 only)
│   └── worker.py                      # Python Worker (part-5 only)
├── typescript-code/                   # TypeScript Activities and Worker implementation for polyglot (part-5 only)
│   ├── activities.js                  # TypeScript Activities (part-5 only)
│   └── worker.js                      # TypeScript Worker (part-5 only)
├── workers/                           # Single worker, polyglot in part-5
│   ├── main.go                        # Go worker (all branches)
│  
├── cmd/
│   ├── starter/main.go                # CLI starter (evolves each branch)
│   └── tools/                         # Demo utilities
│       ├── approve.go
│       ├── status.go
│       └── cleanup.go
└── docs/
    ├── part-1-guide.md                 # Branch-specific instructions
    ├── part-2-guide.md
    ├── part-3-guide.md
    ├── part-4-guide.md
    └── part-5-guide.md
```

### Git Branch Strategy
```
main (setup + foundation)
├── part-1-basic-pipeline
├── part-2-human-in-the-loop
├── part-3-production-features  
├── part-4-crash-resilience (tag pointing to part-1)
└── part-5-polyglot-finale
```

**Branch Progression**:
- **main**: Foundation setup, sample app, basic project structure
- **part-1-basic-pipeline**: Basic Docker pipeline workflow + activities
- **part-2-human-in-the-loop**: Adds Kubernetes deployment + approval workflow
- **part-3-production-features**: Adds timers, scheduling, health checks
- **part-4-crash-resilience**: Git tag pointing to part-1 (for clean crash demo)
- **part-5-polyglot-finale**: Complete rewrite with multi-language workers

---

## Implementation Phases (Branch-Based)

### Phase 1: Foundation (main branch) - 60 minutes
**Goal**: Create the basic infrastructure and sample application that all demos build upon

#### 1.1 Repository Initialization
```bash
# Commands to provide to Claude Code:
mkdir temporal-cicd-workshop && cd temporal-cicd-workshop
git init
```

**Context for Claude Code**:
> Initialize the main branch with the complete workshop foundation. This will be the base that all other branches build upon. Create:
> 1. Temporal CLI provides a development server for running Temporal Workflows
> 2. Kubernetes setup scripts for staging/production namespaces  
> 3. Sample Go application with health endpoints and tests
> 4. Dockerfile and K8s manifests for the sample app
> 5. Shared Go types and utilities that will be used across all demos
> 6. Basic project structure with empty workflow/activity files
> 7. Comprehensive README with branch navigation instructions

#### 1.2 Workshop Navigation Setup
**Context for Claude Code**:
> Create a WORKSHOP_GUIDE.md that explains how to navigate between branches for each demo. Include git commands and what each branch demonstrates. This will be the instructor's guide for delivering the workshop.

### Phase 2: Part 1 Branch (part-1-basic-pipeline) - 90 minutes
**Goal**: Implement basic Docker build pipeline

#### 2.1 Branch from Main
```bash
git checkout -b part-1-basic-pipeline
```

**Context for Claude Code**:
> Starting from the main branch foundation, implement the basic Docker pipeline. Create or modify:
> 1. **workflows/pipeline.go**: Simple sequential workflow (Build -> Test -> Push)
> 2. **activities/docker.go**: Docker build, test, and registry push activities
> 3. **workers/main.go**: Worker that registers Docker activities
> 4. **cmd/starter/main.go**: CLI to trigger basic pipeline workflows
> 5. **docs/part-1-guide.md**: Instructions for this specific demo

> Focus on simplicity and clear demonstration of Temporal basics. Include comprehensive logging and error handling to make the demo engaging.

### Phase 3: Part 2 Branch (part-2-human-in-the-loop) - 60 minutes  
**Goal**: Add Kubernetes deployment with human approval gates

#### 3.1 Branch from Part 1
```bash
git checkout part-1-basic-pipeline
git checkout -b part-2-human-in-the-loop
```

**Context for Claude Code**:
> Building on part-1, add Kubernetes and approval capabilities. Create or modify:
> 1. **workflows/pipeline.go**: Extend to include K8s deployment + approval workflow
> 2. **activities/kubernetes.go**: NEW - K8s deployment and service management
> 3. **activities/approval.go**: NEW - Human-in-the-loop approval activity
> 4. **cmd/starter/main.go**: Add approval commands (-action=approve/reject)
> 5. **docs/part-2-guide.md**: Demo instructions including approval workflow

> Key changes from part-1: workflow now deploys to staging automatically, then waits for human approval before production deployment.

### Phase 4: Part 3 Branch (part-3-production-features) - 45 minutes
**Goal**: Add production-ready features (timers, scheduling, rollbacks)

#### 4.1 Branch from Part 2  
```bash
git checkout part-2-human-in-the-loop
git checkout -b part-3-production-features
```

**Context for Claude Code**:
> Building on part-2, add production-ready operational features. Create or modify:
> 1. **workflows/pipeline.go**: Add deployment windows, health monitoring, rollback timers
> 2. **activities/monitoring.go**: NEW - Health checks and rollback operations
> 3. **cmd/starter/main.go**: Add validation commands (-action=validate)
> 4. **docs/part-3-guide.md**: Instructions for timer and scheduling demos

> Key additions: deployment window checking, automatic rollback timers, manual validation to cancel rollbacks, selector patterns for concurrent operations.

### Phase 5: Part 5 Branch (part-5-polyglot-finale) - 90 minutes
**Goal**: Complete rewrite with multi-language workers

#### 5.1 Branch from Part 3
```bash  
git checkout part-3-production-features
git checkout -b part-5-polyglot-finale
```

**Context for Claude Code**:
> Create a comprehensive polyglot implementation that demonstrates cross-language coordination. This is a significant rewrite:

> 1. **workflows/pipeline.go**: Orchestrates activities across Go, Python, and Node.js workers
> 2. **activities/docker.go** & **activities/kubernetes.go**: Remain Go-based for infrastructure
> 3. **activities/python/test_runner.py**: NEW - Python testing and security activities  
> 4. **activities/python/security_scan.py**: NEW - Container security scanning
> 5. **activities/nodejs/notifications.js**: NEW - Slack/email notifications
> 6. **activities/nodejs/monitoring.js**: NEW - Monitoring setup
> 7. **workers/main.go**: Go worker for infrastructure operations
> 8. **workers/main.py**: NEW - Python worker for testing/security
> 9. **workers/main.js**: NEW - Node.js worker for notifications/monitoring
> 10. **docs/part-5-guide.md**: Instructions for running all three workers

> This demonstrates each language handling its strengths: Go for infrastructure, Python for data/testing, Node.js for async I/O.

### Phase 6: Part 4 Tag (crash resilience)
**Goal**: Clean crash demonstration without code changes

#### 6.1 Create Tag Pointing to Part 1
```bash
git checkout part-1-basic-pipeline  
git tag part-4-crash-resilience
```

**Context for Claude Code**:
> Create documentation explaining that part-4 uses the same code as part-1 but focuses on demonstrating crash resilience. Create docs/part-4-guide.md with specific instructions for the crash demonstration (when to kill the worker, how to restart, what to observe).

---

## Demo-by-Demo Implementation (Branch-Based)

### Part 1: Basic Docker Pipeline (part-1-basic-pipeline branch)
**Implementation Time**: 90 minutes
**Complexity**: Low
**Git Strategy**: Branch from main

#### Goals
- Demonstrate basic Temporal workflow concepts
- Show Docker integration with proper error handling
- Introduce Temporal UI and workflow visibility

#### Key Changes from Main Branch
1. Implement `workflows/pipeline.go` with basic sequential workflow
2. Create `activities/docker.go` with build, test, push activities
3. Enhance `workers/main.go` to register Docker activities
4. Update `cmd/starter/main.go` with pipeline trigger commands

#### Claude Code Instructions
```bash
git checkout main
git checkout -b part-1-basic-pipeline
```

**Context for Claude Code**:
> Building on the main branch foundation, implement a basic Temporal workflow for Docker operations. The workflow should be simple and demonstrate core Temporal concepts clearly.

> **Modify these files**:
> 1. **workflows/pipeline.go**: Create BasicPipelineWorkflow with sequential activities
> 2. **activities/docker.go**: Implement BuildDockerImage, TestDockerContainer, PushToRegistry
> 3. **workers/main.go**: Register Docker activities and start worker
> 4. **cmd/starter/main.go**: Add commands to trigger pipeline workflows
> 5. **docs/part-1-guide.md**: Instructions for running this demo

> **Key requirements**:
> - Simple sequential execution: Build -> Test -> Push
> - Comprehensive logging for demo visibility
> - Activity retry policies for Docker operations
> - Idempotent activities (safe to retry)
> - CLI commands: `go run cmd/starter/main.go -action=create -image=demo-app:v1.0.0`

#### Expected Demo Flow
1. Start worker: `go run workers/main.go`
2. Trigger pipeline: `go run cmd/starter/main.go -action=create -image=demo-app:v1.0`
3. Show Temporal UI with workflow progress
4. Point out retry policies and activity details

---

### Part 2: Human Approval Integration (part-2-human-in-the-loop branch)
**Implementation Time**: 60 minutes
**Complexity**: Medium
**Git Strategy**: Branch from part-1-basic-pipeline

#### Goals
- Demonstrate human-in-the-loop workflows
- Show Kubernetes integration
- Introduce signals and queries

#### Key Changes from Part 1 Branch
1. Extend workflow to include Kubernetes deployment
2. Add approval activities and signal handling
3. Enhance starter with approval commands
4. Add Kubernetes deployment activities

#### Claude Code Instructions
```bash
git checkout part-1-basic-pipeline
git checkout -b part-2-human-in-the-loop
```

**Context for Claude Code**:
> Building on part-1, extend the pipeline to include Kubernetes deployment with human approval gates. This demonstrates how Temporal handles long-running processes with human interaction.

> **Create these new files**:
> 1. **activities/kubernetes.go**: K8s deployment, status checking, service discovery
> 2. **activities/approval.go**: Human approval activity using signals

> **Modify these existing files**:
> 1. **workflows/pipeline.go**: Extend with K8s deployment + approval logic
> 2. **cmd/starter/main.go**: Add approval commands (-action=approve/reject/status)
> 3. **docs/part-2-guide.md**: Demo instructions with approval workflow

> **Key requirements**:
> - Automatic staging deployment after successful build/test
> - Human approval gate before production deployment
> - Workflow waits indefinitely for approval signal
> - CLI approval: `go run cmd/starter/main.go -action=approve -workflow=<id>`
> - Include staging URL in approval decision

#### Expected Demo Flow
1. Start enhanced workflow with production target
2. Show automatic staging deployment
3. Demonstrate approval CLI and workflow resumption
4. Show workflow state preservation during approval wait

---

### Part 3: Production Features (part-3-production-features branch)
**Implementation Time**: 45 minutes
**Complexity**: Medium-High
**Git Strategy**: Branch from part-2-human-in-the-loop

#### Goals
- Demonstrate durable timers
- Show advanced workflow patterns (selectors)
- Include production-ready features

#### Key Changes from Part 2 Branch
1. Add deployment window checking
2. Implement rollback timers and health monitoring
3. Use selector patterns for concurrent operations
4. Add manual validation to cancel rollbacks

#### Claude Code Instructions
```bash
git checkout part-2-human-in-the-loop
git checkout -b part-3-production-features
```

**Context for Claude Code**:
> Building on part-2, add production-ready operational features that demonstrate Temporal's advanced capabilities for real-world deployments.

> **Create these new files**:
> 1. **activities/monitoring.go**: Health checks, rollback operations, deployment windows

> **Modify these existing files**:
> 1. **workflows/pipeline.go**: Add deployment windows, health monitoring, rollback timers
> 2. **cmd/starter/main.go**: Add validation commands (-action=validate)
> 3. **docs/part-3-guide.md**: Instructions for timer and scheduling demos

> **Key requirements**:
> - Deployment window checking (sleep until valid window)
> - Automatic rollback timer (30 min, demo with 30 seconds)
> - Manual validation to cancel rollback timer
> - Health monitoring after production deployment
> - Selector pattern for handling timer vs validation signal
> - CLI validation: `go run cmd/starter/main.go -action=validate -workflow=<id>`

#### Expected Demo Flow
1. Show deployment window checking (configure for demo timing)
2. Deploy to production with rollback timer
3. Demonstrate manual validation canceling rollback
4. Show timer functionality in Temporal UI

---

### Part 4: Crash Resilience (part-4-crash-resilience tag)
**Implementation Time**: 5 minutes (documentation only)
**Complexity**: Low (demo execution)
**Git Strategy**: Tag pointing to part-1-basic-pipeline

#### Goals
- Demonstrate Temporal's core value proposition
- Show workflow durability and deterministic replay
- Create memorable "wow" moment for attendees

#### Implementation
```bash
git checkout part-1-basic-pipeline
git tag part-4-crash-resilience
```

**Context for Claude Code**:
> Create documentation for the crash resilience demonstration. This uses the same code as part-1 but focuses on showing how Temporal handles worker crashes.

> **Create**:
> 1. **docs/part-4-guide.md**: Detailed crash demonstration instructions

> **Key content**:
> - When to kill the worker process during demo
> - How to restart the worker cleanly
> - What to observe in Temporal UI
> - Talking points about deterministic replay
> - Commands for clean crash demo execution

#### Expected Demo Flow
1. Start workflow: `go run cmd/starter/main.go -action=create -image=demo-app:crash-test`
2. Kill worker during Docker build: `pkill -f worker`
3. Restart worker: `go run workers/main.go`
4. Show workflow continues exactly where it left off
5. Highlight no duplicate Docker operations occurred

---

### Part 5: Polyglot Coordination (part-5-polyglot-finale branch)
**Implementation Time**: 90 minutes
**Complexity**: High
**Git Strategy**: Branch from part-3-production-features

#### Goals
- Demonstrate cross-language activity execution
- Show specialized workers for different concerns
- Provide compelling finale showcasing Temporal's flexibility

#### Key Changes from Part 3 Branch
1. Complete rewrite of worker architecture
2. Add Python worker for testing and security
3. Add Node.js worker for notifications and monitoring
4. Implement cross-language activity coordination

#### Claude Code Instructions
```bash
git checkout part-3-production-features
git checkout -b part-5-polyglot-finale
```

**Context for Claude Code**:
> Create a comprehensive polyglot implementation that demonstrates Temporal's ability to coordinate activities across different programming languages and technology stacks.

> **Create these new files**:
> 1. **activities/python/test_runner.py**: Integration testing with pytest
> 2. **activities/python/security_scan.py**: Container security scanning
> 3. **activities/nodejs/notifications.js**: Slack/email notifications
> 4. **activities/nodejs/monitoring.js**: Monitoring and alerting setup
> 5. **workers/main.py**: Python worker for testing/security activities
> 6. **workers/main.js**: Node.js worker for notifications/monitoring

> **Modify these existing files**:
> 1. **workflows/pipeline.go**: Orchestrate activities across all three workers
> 2. **workers/main.go**: Focus on Docker and Kubernetes activities only
> 3. **cmd/starter/main.go**: Handle coordination of multiple workers
> 4. **docs/part-5-guide.md**: Instructions for running all three workers

> **Worker Specialization**:
> - **Go Worker**: Infrastructure operations (Docker, Kubernetes)
> - **Python Worker**: Testing, security scanning, data analysis
> - **Node.js Worker**: Notifications, monitoring, async I/O operations

> **Key requirements**:
> - Single workflow coordinates all three language workers
> - Each worker registers only language-appropriate activities
> - Graceful error handling across workers
> - Rich logging showing which worker executes each activity
> - Startup instructions for all three workers in separate terminals

#### Expected Demo Flow
1. Start all three workers in separate terminals
2. Trigger polyglot workflow
3. Show Temporal UI with activities from different workers
4. Highlight language-specific strengths and specialization
5. Demonstrate failure handling across workers

---

## Supporting Infrastructure

### Temporal Server Setup
Use Temporal CLI and the command `temporal server start-dev --db-file temporal.db

### Kubernetes Setup
**File**: `setup/setup-k8s.sh`

**Context for Claude Code**:
> Create a setup script that prepares Kubernetes environment for the workshop. Create staging and production namespaces, set up service accounts if needed, and configure any necessary RBAC. Should work with Docker Desktop Kubernetes or kind clusters.

### Demo Utilities
**Files**: `tools/approve.go`, `tools/status.go`, `tools/cleanup.go`

**Context for Claude Code**:
> Create utility tools for workshop demonstration:
> 1. **approve.go**: CLI tool for approving/rejecting deployments
> 2. **status.go**: Shows current status of workflows and deployments
> 3. **cleanup.go**: Cleans up Docker images, Kubernetes deployments, and Temporal workflows between demos

---

## Testing & Validation

### Pre-Workshop Testing Checklist
1. **Environment Validation**
   - [ ] Temporal server starts correctly
   - [ ] Kubernetes namespaces created
   - [ ] Docker daemon accessible
   - [ ] Sample app builds and runs

2. **Part 1 Testing**
   - [ ] Docker build succeeds
   - [ ] Container tests pass
   - [ ] Registry push works
   - [ ] Retry policies function correctly

3. **Part 2 Testing**
   - [ ] Kubernetes deployment succeeds
   - [ ] Approval workflow pauses correctly
   - [ ] Approval/rejection both work
   - [ ] Service URLs accessible

4. **Part 3 Testing**
   - [ ] Deployment window logic works
   - [ ] Rollback timer functions
   - [ ] Manual validation cancels timer
   - [ ] Health checks operate correctly

5. **Part 4 Testing**
   - [ ] Worker crash during build
   - [ ] Workflow resumes correctly
   - [ ] No duplicate operations occur

6. **Part 5 Testing**
   - [ ] All three workers start
   - [ ] Cross-language coordination works
   - [ ] Failure handling across workers
   - [ ] Activities assigned to correct workers

### Validation Scripts
**Context for Claude Code**:
> Create validation scripts that can test each demo independently:
> 1. **validate-part-1.sh**: Runs part-1 end-to-end and validates results
> 2. **validate-part-2.sh**: Tests approval workflow with automated approval
> 3. **validate-part-3.sh**: Tests timer functionality with fast timers
> 4. **validate-part-5.sh**: Starts all workers and validates polyglot coordination

---

## Workshop Delivery Guide

### Pre-Workshop Setup (15 minutes before)
1. **Start Infrastructure**
   ```bash
   cd setup
   ./setup-k8s.sh
   ```

2. **Validate Environment**
   ```bash
   ./validate-environment.sh
   ```

3. **Prepare Demo State**
   ```bash
   cd sample-app
   docker build -t demo-app:base .
   ```

### Part 1: Basic Pipeline (0-15 minutes)
**Key Talking Points**:
- Temporal workflows as code
- Activity retry policies
- Workflow visibility in UI

**Demo Commands**:
```bash
# Terminal 1: Start worker
go run worker/main.go

# Terminal 2: Trigger workflow
go run starter/main.go -image=demo-app:v1.0.0
```

**What to Show**:
- Workflow execution in Temporal UI
- Activity retry on simulated failure
- Docker images being created
- Clean sequential execution

### Part 2: Human Approval (15-30 minutes)
**Key Talking Points**:
- Human-in-the-loop workflows
- Long-running processes
- Signals and queries

**Demo Commands**:
```bash
# Start enhanced worker
go run worker/main.go

# Trigger with approval
go run starter/main.go -image=demo-app:v2.0.0 -env=production

# Show staging deployment
kubectl get pods -n staging

# Approve deployment
go run starter/main.go -action=approve -workflow=<id>
```

**What to Show**:
- Workflow paused on approval
- Staging service running
- Approval process
- Production deployment proceeding

### Part 3: Production Features (30-45 minutes)
**Key Talking Points**:
- Durable timers
- Deployment windows
- Automatic rollback capabilities

**Demo Commands**:
```bash
# Start timer-enabled worker
cd part-3 && go run worker/main.go

# Trigger with rollback timer
go run starter/main.go -image=demo-app:v3.0.0 -env=production

# Show rollback timer in UI
# Validate deployment to cancel timer
go run starter/main.go -action=validate -workflow=<id>
```

**What to Show**:
- Deployment window checking
- Rollback timer in UI
- Manual validation canceling timer
- Health monitoring

### Part 4: Crash Resilience (45-60 minutes)
**Key Talking Points**:
- Temporal's core value proposition
- Deterministic replay
- No lost work

**Demo Commands**:
```bash
# Start workflow
cd part-4 && go run worker/main.go &
WORKER_PID=$!

# Trigger build
go run starter/main.go -image=demo-app:crash-test

# Kill worker mid-build
kill $WORKER_PID

# Restart worker
go run worker/main.go
```

**What to Show**:
- Workflow continuing from exact point
- No duplicate Docker operations
- State preservation across crashes

### Part 5: Polyglot Finale (60-75 minutes)
**Key Talking Points**:
- Cross-language coordination
- Specialized workers
- Technology diversity

**Demo Commands**:
```bash
# Start all workers (3 terminals)
cd part-5
go run workers/go-worker/main.go &
python w &
node workers/nodejs-worker/main.js &

# Trigger polyglot workflow
go run starter/main.go -image=demo-app:polyglot
```

**What to Show**:
- Activities distributed across workers
- Language-specific strengths
- Unified coordination
- Rich logging showing worker assignments

### Q&A and Wrap-up (75-90 minutes)
**Key Discussion Points**:
- Real-world applications
- Integration patterns
- Getting started resources
- Common pitfalls and solutions

---

## Success Metrics

### Technical Success
- [ ] All demos execute without errors
- [ ] Temporal UI shows workflow progress clearly
- [ ] Docker images build and deploy successfully
- [ ] Kubernetes deployments reach ready state
- [ ] Cross-language coordination works seamlessly

### Educational Success
- [ ] Attendees understand Temporal's value proposition
- [ ] Questions indicate comprehension of concepts
- [ ] Interest in implementing Temporal for their use cases
- [ ] Requests for additional resources/follow-up

### Delivery Success
- [ ] Workshop stays on schedule
- [ ] All key features demonstrated
- [ ] Engaging and interactive presentation
- [ ] Repository provides immediate value to attendees

This implementation plan provides a comprehensive roadmap for building a compelling Temporal workshop that demonstrates real-world value through practical CI/CD scenarios. Each phase builds upon the previous, culminating in a sophisticated demonstration of Temporal's capabilities for platform engineering teams.