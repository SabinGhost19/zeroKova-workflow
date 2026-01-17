# Flux CD Deployment Configuration

This directory contains Flux CD configuration for GitOps-based deployment of test-workflow.

## Directory Structure

```
flux/
├── base/                      # base flux resources
│   ├── kustomization.yaml     # kustomize aggregation
│   ├── namespace.yaml         # target namespace
│   ├── source.yaml            # git repository source
│   └── helmrelease.yaml       # helm release definition
└── overlays/                  # environment-specific overrides
    ├── development/           # development environment
    │   └── kustomization.yaml
    ├── staging/               # staging environment
    │   └── kustomization.yaml
    └── production/            # production environment
        └── kustomization.yaml
```

## Prerequisites

1. Kubernetes cluster with Flux CD installed
2. Flux CLI (`flux`) installed locally
3. Access to the git repository

## Installation

### Bootstrap Flux (if not already installed)

```bash
flux bootstrap github \
  --owner=sabinghosty19 \
  --repository=test-workflow \
  --path=flux/overlays/development \
  --personal
```

### Deploy to Development

```bash
flux create kustomization test-workflow-dev \
  --source=GitRepository/flux-system \
  --path="./flux/overlays/development" \
  --prune=true \
  --interval=5m
```

### Deploy to Staging

```bash
flux create kustomization test-workflow-staging \
  --source=GitRepository/flux-system \
  --path="./flux/overlays/staging" \
  --prune=true \
  --interval=5m
```

### Deploy to Production

```bash
flux create kustomization test-workflow-prod \
  --source=GitRepository/flux-system \
  --path="./flux/overlays/production" \
  --prune=true \
  --interval=5m
```

## Monitoring Deployments

```bash
# check flux resources
flux get all

# check helmrelease status
flux get helmreleases -n flux-system

# check kustomization status
flux get kustomizations

# watch reconciliation
flux logs --follow
```

## Troubleshooting

```bash
# reconcile immediately
flux reconcile kustomization test-workflow-dev

# check events
kubectl events -n test-workflow

# describe helmrelease
kubectl describe helmrelease test-workflow -n flux-system
```
