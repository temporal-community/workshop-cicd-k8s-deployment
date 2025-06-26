# Python Kubernetes Activities

This directory contains the Python implementation of Kubernetes activities for the Temporal CI/CD workshop. It demonstrates how to implement the same functionality as the Go Kubernetes activities using Python and the Temporal Python SDK.

## Overview

The Python worker implements four Kubernetes activities:
- `deploy_to_kubernetes`: Deploy applications to staging/production namespaces
- `check_deployment_status`: Query deployment status via kubectl
- `rollback_deployment`: Perform rollback operations with automatic fallback
- `get_service_url`: Retrieve service URLs for deployed applications

## Setup

### Prerequisites
- Python 3.8 or higher
- [uv](https://github.com/astral-sh/uv) package manager
- kubectl configured for your cluster
- Temporal server running (default: localhost:7233)

### Installation with uv

1. Install uv if you haven't already:
```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

2. Install dependencies:
```bash
uv sync
```

## Running the Worker

### Using uv
```bash
uv run python src/worker.py
```

### Using traditional Python (after installing dependencies)
```bash
python src/worker.py
```

### Environment Variables
- `TEMPORAL_HOST`: Temporal server address (default: localhost:7233)

## Project Structure

```
python-activities/
├── pyproject.toml        # Python project configuration and dependencies
├── src/
│   ├── __init__.py
│   ├── activities/
│   │   ├── __init__.py
│   │   └── kubernetes.py # Kubernetes activity implementations
│   ├── types.py          # Python dataclass type definitions
│   └── worker.py         # Worker startup and configuration
└── README.md
```

## Activities

### deploy_to_kubernetes
Deploys applications to Kubernetes with the following steps:
1. Update deployment with new image (or create if doesn't exist)
2. Wait for rollout to complete with timeout
3. Ensure service exists (create if needed)
4. Retrieve service URL (LoadBalancer or NodePort fallback)
5. Verify deployment health

### check_deployment_status
Queries deployment status including replica counts and readiness state.

### rollback_deployment
Performs deployment rollbacks with these steps:
1. Check if deployment exists
2. Execute `kubectl rollout undo`
3. Wait for rollback completion
4. Verify rollback success
5. Fallback to deletion if rollback fails

### get_service_url
Retrieves service URLs with environment-appropriate protocols (HTTP for staging, HTTPS for production).

## Kubernetes Integration

The activities use `kubectl` commands via Python's `asyncio.subprocess` for Kubernetes operations:
- Deployment management (create, update, rollback)
- Service management (create, query)
- Status monitoring and health checks
- URL retrieval with LoadBalancer/NodePort fallback

## Integration with Go Workflow

This Python worker connects to the same Temporal cluster and task queue (`cicd-task-queue`) as the Go workers. The Go workflow can call these Python activities seamlessly, demonstrating Temporal's polyglot capabilities.

## Development

### Code Formatting
```bash
uv run black src/
```

### Type Checking
```bash
uv run mypy src/
```

### Testing
```bash
uv run pytest
```

## Error Handling

The activities implement comprehensive error handling:
- Kubectl command failures with detailed logging
- Deployment rollback fallbacks (deletion when undo fails)
- Service URL fallbacks (NodePort when LoadBalancer unavailable)
- Heartbeat recording for long-running operations
- Structured logging for observability

## Notes

- All activities are async and use Python's asyncio for concurrent operations
- Kubectl operations include proper timeout handling and error recovery
- Activities record heartbeats for Temporal's progress tracking
- Service URL generation adapts to different Kubernetes environments
- The worker handles graceful shutdown on interrupt signals