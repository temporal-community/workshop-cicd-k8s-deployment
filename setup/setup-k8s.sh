#!/bin/bash

echo "Setting up Kubernetes namespaces for Temporal CI/CD Workshop..."

# Create namespaces
kubectl create namespace staging --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace production --dry-run=client -o yaml | kubectl apply -f -

# Label namespaces
kubectl label namespace staging environment=staging --overwrite
kubectl label namespace production environment=production --overwrite

echo "✅ Kubernetes namespaces created successfully!"
echo ""
echo "Namespaces:"
kubectl get namespaces | grep -E "(staging|production)"