#!/bin/bash

echo "Cleaning up workshop resources..."

# Delete Kubernetes deployments and services
echo "Cleaning up Kubernetes resources..."
kubectl delete deployment,service -n staging --all 2>/dev/null
kubectl delete deployment,service -n production --all 2>/dev/null

# Remove Docker images created during workshop
echo "Cleaning up Docker images..."
docker images | grep "demo-app" | awk '{print $3}' | xargs -r docker rmi -f 2>/dev/null

# Clean up Temporal workflows (requires temporal CLI)
echo "Note: To clean up Temporal workflows, run:"
echo "  temporal workflow terminate --query 'WorkflowType=\"CICDPipelineWorkflow\"'"

echo "✅ Cleanup completed!"