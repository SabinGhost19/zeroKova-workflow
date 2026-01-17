#!/bin/bash

# Deploy all resources to Kubernetes

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
K8S_DIR="$PROJECT_ROOT/k8s"

echo "Deploying test-workflow to Kubernetes..."

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if kustomize is available
if ! command -v kustomize &> /dev/null; then
    echo "Using kubectl kustomize..."
    kubectl apply -k "$K8S_DIR/base"
else
    echo "Using kustomize..."
    kustomize build "$K8S_DIR/base" | kubectl apply -f -
fi

echo ""
echo "Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app.kubernetes.io/part-of=test-workflow -n test-workflow --timeout=300s || true

echo ""
echo "Deployment status:"
kubectl get pods -n test-workflow

echo ""
echo "Services:"
kubectl get svc -n test-workflow

echo ""
echo "Ingress:"
kubectl get ingress -n test-workflow

echo ""
echo "Deployment complete!"
echo "Add '127.0.0.1 test-workflow.local' to /etc/hosts to access the application"
