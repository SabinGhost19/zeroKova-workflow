# GitHub Actions CI/CD Pipeline Guide

This document provides comprehensive documentation for the CI/CD pipeline configuration used in the test-workflow project. The pipeline automates testing, building, security scanning, and releasing of all microservices.

---

## Table of Contents

1. [Pipeline Architecture](#pipeline-architecture)
2. [Workflow Descriptions](#workflow-descriptions)
3. [Repository Configuration](#repository-configuration)
4. [Docker Image Management](#docker-image-management)
5. [Security Scanning](#security-scanning)
6. [Release Automation](#release-automation)
7. [Dependency Management](#dependency-management)
8. [Troubleshooting](#troubleshooting)
9. [Command Reference](#command-reference)
10. [References](#references)

---

## Pipeline Architecture

### Workflow Structure

```
.github/
├── README.md                    # This document
├── docs/
│   ├── RENOVATE.md              # Renovate bot configuration guide
│   ├── SEMANTIC_RELEASE.md      # Semantic versioning documentation
│   └── TRIVY.md                 # Security scanning guide
└── workflows/
    ├── ci-test.yaml             # Multi-language test runner
    ├── ci-build.yaml            # Docker image builder
    ├── ci-security.yaml         # Trivy vulnerability scanner
    ├── ci-helm.yaml             # Helm chart validator
    ├── release.yaml             # Semantic release automation
    ├── renovate.yaml            # Dependency update bot
    └── pipeline-summary.yaml    # Status aggregation
```

### Pipeline Flow Diagram

```
                            CI/CD Pipeline Architecture
                            ===========================

    +-------------------+
    |   Developer       |
    |   Push/PR         |
    +--------+----------+
             |
             v
    +--------+----------+     +-------------------+
    |                   |     |                   |
    |   ci-test.yaml    +---->+   Test Results    |
    |   (all services)  |     |   (per language)  |
    |                   |     |                   |
    +--------+----------+     +-------------------+
             |
             | on: push to main
             v
    +--------+----------+     +-------------------+
    |                   |     |                   |
    |   ci-build.yaml   +---->+   GHCR Images     |
    |   (Docker build)  |     |   (tagged)        |
    |                   |     |                   |
    +--------+----------+     +-------------------+
             |
             | workflow_run: completed
             v
    +--------+----------+     +-------------------+
    |                   |     |                   |
    |  ci-security.yaml +---->+   Trivy Reports   |
    |  (vulnerability)  |     |   (SARIF/HTML)    |
    |                   |     |                   |
    +--------+----------+     +-------------------+
             |
             v
    +--------+----------+     +-------------------+
    |                   |     |                   |
    |   release.yaml    +---->+   GitHub Release  |
    |   (semantic-rel)  |     |   + CHANGELOG     |
    |                   |     |                   |
    +-------------------+     +-------------------+
```

### Test Matrix Structure

```
                        Test Workflow Matrix
                        ====================

    +------------------------------------------------------------------+
    |                        ci-test.yaml                              |
    |                                                                  |
    |   +------------------+  +------------------+  +----------------+ |
    |   |                  |  |                  |  |                | |
    |   |  API Gateway     |  |  Order Service   |  |  Inventory     | |
    |   |  (Go 1.21)       |  |  (Java 21)       |  |  Service       | |
    |   |                  |  |                  |  |  (C# .NET 8)   | |
    |   |  go test ./...   |  |  mvn test        |  |  dotnet test   | |
    |   |                  |  |                  |  |                | |
    |   +------------------+  +------------------+  +----------------+ |
    |                                                                  |
    |   +------------------+  +------------------+                     |
    |   |                  |  |                  |                     |
    |   |  Notification    |  |  Frontend        |                     |
    |   |  Service         |  |  (Node 20)       |                     |
    |   |  (Go 1.21)       |  |                  |                     |
    |   |  go test ./...   |  |  npm test        |                     |
    |   |                  |  |                  |                     |
    |   +------------------+  +------------------+                     |
    |                                                                  |
    +------------------------------------------------------------------+
```

### Build and Push Flow

```
                        Build Workflow Detail
                        =====================

    +-------------+     +------------------+     +------------------+
    |             |     |                  |     |                  |
    |  Checkout   +---->+  Detect Changes  +---->+  Matrix Build    |
    |  Code       |     |  (paths filter)  |     |  (parallel)      |
    |             |     |                  |     |                  |
    +-------------+     +------------------+     +--------+---------+
                                                          |
                                                          v
    +------------------+     +------------------+     +---+------------+
    |                  |     |                  |     |                |
    |  Push to GHCR    +<----+  Tag Images      +<----+  Docker Build  |
    |                  |     |  (sha, latest)   |     |  (multi-stage) |
    |                  |     |                  |     |                |
    +------------------+     +------------------+     +----------------+
```

---

## Workflow Descriptions

### CI - Tests (ci-test.yaml)

Runs automated tests for all microservices using language-specific test frameworks.

| Trigger | Condition |
|---------|-----------|
| push | main branch, service paths |
| pull_request | main branch, service paths |
| workflow_dispatch | Manual trigger |

**Test Commands by Service:**

| Service | Language | Test Command |
|---------|----------|--------------|
| api-gateway | Go | `go test -v -race ./...` |
| order-service | Java | `mvn test -B` |
| inventory-service | C# | `dotnet test --verbosity normal` |
| notification-service | Go | `go test -v ./...` |
| frontend | Node.js | `npm test -- --coverage` |

### CI - Build and Push (ci-build.yaml)

Builds Docker images and pushes them to GitHub Container Registry.

| Trigger | Condition |
|---------|-----------|
| push | main branch only |
| workflow_dispatch | Manual with force_build_all option |

**Build Matrix:**

```yaml
strategy:
  matrix:
    service:
      - api-gateway
      - order-service
      - inventory-service
      - notification-service
      - frontend
```

### CI - Security Scan (ci-security.yaml)

Runs Trivy vulnerability scanner against all container images and infrastructure code.

| Trigger | Condition |
|---------|-----------|
| workflow_run | After ci-build.yaml completes |
| workflow_dispatch | Manual trigger |

**Scan Targets:**

| Target | Type | Output |
|--------|------|--------|
| Docker images | Container | SARIF, HTML, JSON |
| Helm charts | IaC | SARIF |
| Dockerfiles | Config | SARIF |

### CI - Helm Validation (ci-helm.yaml)

Validates Helm chart syntax and templates.

| Trigger | Condition |
|---------|-----------|
| push | main branch, helm/ path |
| pull_request | main branch, helm/ path |
| workflow_dispatch | Manual trigger |

**Validation Steps:**

1. Lint chart with `helm lint`
2. Template rendering with `helm template`
3. Kubernetes manifest validation with `kubeval`

### Release (release.yaml)

Automates semantic versioning and release creation based on conventional commits.

| Trigger | Condition |
|---------|-----------|
| push | main branch |
| workflow_dispatch | Manual with dry_run option |

### Renovate (renovate.yaml)

Manages automated dependency updates across all services.

| Trigger | Condition |
|---------|-----------|
| schedule | Daily at 06:00 UTC |
| workflow_dispatch | Manual trigger |

---

## Repository Configuration

### Required Settings

Navigate to **Settings** in your GitHub repository and configure the following:

#### Actions Permissions

Path: Settings > Actions > General

| Setting | Value |
|---------|-------|
| Actions permissions | Allow all actions and reusable workflows |
| Workflow permissions | Read and write permissions |
| Allow GitHub Actions to create and approve pull requests | Enabled |

#### Secrets Configuration

The pipeline uses automatically provided secrets:

| Secret | Source | Usage |
|--------|--------|-------|
| GITHUB_TOKEN | Automatic | GHCR authentication, API calls, releases |

No additional secrets are required for basic operation.

#### Package Registry

Path: Settings > Packages

| Setting | Value |
|---------|-------|
| Package visibility | Inherit access from source repository |

### Branch Protection Rules

Recommended configuration for the main branch:

Path: Settings > Branches > Add rule

| Rule | Recommended Value |
|------|-------------------|
| Branch name pattern | main |
| Require pull request before merging | Enabled |
| Require status checks to pass | Enabled |
| Required checks | Test api-gateway, Test order-service, Test inventory-service, Test notification-service, Test frontend |
| Require branches to be up to date | Enabled |
| Include administrators | Enabled |

---

## Docker Image Management

### Registry Location

All images are published to GitHub Container Registry (GHCR):

```
ghcr.io/{owner}/test-workflow/{service}:{tag}
```

### Tag Naming Convention

| Tag Pattern | Description | Example |
|-------------|-------------|---------|
| latest | Most recent main branch build | ghcr.io/.../api-gateway:latest |
| sha-{hash} | Specific commit SHA (first 7 chars) | ghcr.io/.../api-gateway:sha-abc1234 |
| v{x.y.z} | Semantic version from release | ghcr.io/.../api-gateway:v1.2.3 |
| {branch} | Branch name (non-main) | ghcr.io/.../api-gateway:feature-auth |

### Pulling Images

```bash
# Latest version
docker pull ghcr.io/sabinghosty19/test-workflow/api-gateway:latest

# Specific version
docker pull ghcr.io/sabinghosty19/test-workflow/api-gateway:v1.2.3

# Specific commit
docker pull ghcr.io/sabinghosty19/test-workflow/api-gateway:sha-abc1234

# All services (latest)
for svc in api-gateway order-service inventory-service notification-service frontend; do
  docker pull ghcr.io/sabinghosty19/test-workflow/$svc:latest
done
```

### Image Size Optimization

Each service uses multi-stage Docker builds to minimize image size:

| Service | Base Image | Final Size (approx) |
|---------|------------|---------------------|
| api-gateway | golang:alpine / scratch | 15 MB |
| order-service | eclipse-temurin:21-jre | 250 MB |
| inventory-service | mcr.microsoft.com/dotnet/aspnet:8.0 | 200 MB |
| notification-service | golang:alpine / alpine | 20 MB |
| frontend | nginx:alpine | 25 MB |

---

## Security Scanning

### Trivy Scanner Configuration

The security workflow runs Trivy against multiple targets:

```
                    Security Scan Flow
                    ==================

    +------------------+
    |                  |
    |  Docker Images   +------+
    |  (5 services)    |      |
    |                  |      |
    +------------------+      |      +------------------+
                              +----->+                  |
    +------------------+      |      |  Trivy Scanner   |
    |                  |      |      |                  |
    |  Helm Charts     +------+      +--------+---------+
    |  (IaC scan)      |      |               |
    |                  |      |               v
    +------------------+      |      +--------+---------+
                              |      |                  |
    +------------------+      |      |  Report Output   |
    |                  |      |      |  - SARIF         |
    |  Dockerfiles     +------+      |  - HTML          |
    |  (config scan)   |             |  - JSON          |
    |                  |             |                  |
    +------------------+             +------------------+
```

### Report Access

#### GitHub Actions Artifacts

After each scan, reports are available in:

1. Navigate to **Actions** > Select workflow run
2. Scroll to **Artifacts** section
3. Download `trivy-security-reports-all` (combined report with index.html)

#### GitHub Security Tab

SARIF reports are uploaded to GitHub Security:

1. Navigate to **Security** > **Code scanning alerts**
2. View all vulnerabilities across services
3. Filter by severity, service, or vulnerability type

### Report Formats

| Format | Purpose | Location |
|--------|---------|----------|
| HTML | Human-readable report with styling | Artifacts download |
| SARIF | GitHub Security integration | Security tab |
| JSON | Machine processing, CI integration | Artifacts download |

### Severity Levels

| Level | Action |
|-------|--------|
| CRITICAL | Fails workflow (configurable) |
| HIGH | Warning, documented in report |
| MEDIUM | Documented in report |
| LOW | Documented in report |

---

## Release Automation

### Semantic Versioning

The project follows Semantic Versioning (SemVer) with automated releases based on commit messages.

### Conventional Commits Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Version Bump Rules

| Commit Type | Version Bump | Example Commit |
|-------------|--------------|----------------|
| fix: | Patch (0.0.X) | fix: resolve null pointer in order service |
| feat: | Minor (0.X.0) | feat: add bulk order creation endpoint |
| feat!: | Major (X.0.0) | feat!: redesign order API response format |
| BREAKING CHANGE: | Major (X.0.0) | feat: update auth\n\nBREAKING CHANGE: token format changed |

### Release Workflow

```
                        Release Process
                        ===============

    +------------------+     +------------------+     +------------------+
    |                  |     |                  |     |                  |
    |  Analyze Commits +---->+  Calculate Next  +---->+  Update         |
    |  (since last tag)|     |  Version         |     |  CHANGELOG.md   |
    |                  |     |                  |     |                  |
    +------------------+     +------------------+     +--------+---------+
                                                               |
                                                               v
    +------------------+     +------------------+     +--------+---------+
    |                  |     |                  |     |                  |
    |  Re-tag Docker   +<----+  Create GitHub   +<----+  Create Git Tag  |
    |  Images          |     |  Release         |     |  (v1.2.3)        |
    |                  |     |                  |     |                  |
    +------------------+     +------------------+     +------------------+
```

### Release Outputs

Each release creates:

1. **Git Tag**: Annotated tag (e.g., v1.2.3)
2. **GitHub Release**: With auto-generated release notes
3. **CHANGELOG.md**: Updated with new version section
4. **Docker Image Tags**: Images re-tagged with version

---

## Dependency Management

### Renovate Bot Configuration

Renovate automatically manages dependency updates across all services.

### Schedule

| Parameter | Value |
|-----------|-------|
| Schedule | Daily at 06:00 UTC (weekdays) |
| Timezone | UTC |
| Automerge | Enabled for patch and minor updates |

### Update Grouping

Updates are grouped to reduce PR noise:

| Group | Includes | Automerge |
|-------|----------|-----------|
| Go Dependencies | go.mod updates | Patch/Minor |
| Java Dependencies | pom.xml updates | Patch/Minor |
| .NET Dependencies | .csproj updates | Patch/Minor |
| Node Dependencies | package.json updates | Patch/Minor |
| Docker Base Images | Dockerfile FROM updates | Patch only |
| GitHub Actions | workflow action versions | Patch/Minor |

### PR Labels

| Label | Meaning |
|-------|---------|
| dependencies | All dependency update PRs |
| renovate | Created by Renovate bot |
| major-update | Major version bump (requires review) |
| minor-update | Minor version bump |
| patch-update | Patch version bump |
| security | Security-related update |

### Managing Updates

To skip a specific dependency update, comment on the PR:

```
@renovate ignore this dependency
```

To ignore a major version:

```
@renovate ignore this major version
```

---

## Troubleshooting

### Test Failures

```
                    Test Failure Diagnosis
                    ======================

    +------------------+
    |                  |
    |  Check Actions   +---->  View test job logs
    |  Workflow Run    |
    |                  |
    +--------+---------+
             |
             v
    +--------+---------+
    |                  |
    |  Identify        +---->  Note failing test name
    |  Failed Test     |
    |                  |
    +--------+---------+
             |
             v
    +--------+---------+
    |                  |
    |  Run Locally     +---->  Reproduce and fix
    |                  |
    +------------------+
```

**Local Test Commands:**

```bash
# Go services
cd services/api-gateway && go test -v ./...
cd services/notification-service && go test -v ./...

# Java service
cd services/order-service && mvn test

# C# service
cd services/inventory-service && dotnet test

# Frontend
cd frontend && npm test
```

### Build Failures

Common causes and solutions:

| Issue | Cause | Solution |
|-------|-------|----------|
| Dockerfile syntax error | Invalid instruction | Validate with `docker build --check` |
| Missing dependencies | Not in go.mod/package.json | Run dependency install |
| Build timeout | Large image or slow network | Increase timeout or optimize Dockerfile |
| GHCR push failed | Authentication issue | Check GITHUB_TOKEN permissions |

**Local Build Test:**

```bash
# Build single service
docker build -t test-build services/api-gateway

# Build with BuildKit (recommended)
DOCKER_BUILDKIT=1 docker build -t test-build services/api-gateway
```

### Security Scan Issues

| Issue | Solution |
|-------|----------|
| Critical vulnerability found | Update base image or dependency |
| False positive | Add to .trivyignore file |
| Scan timeout | Increase workflow timeout |

### Release Not Triggering

Checklist:

1. Verify commits follow Conventional Commits format
2. Check commit messages do not contain `[skip ci]`
3. Confirm push is to main branch
4. Review release workflow logs for errors

---

## Command Reference

### Manual Workflow Triggers

```bash
# Run tests for all services
gh workflow run ci-test.yaml

# Build and push all images (force rebuild)
gh workflow run ci-build.yaml -f force_build_all=true

# Run security scan
gh workflow run ci-security.yaml

# Validate Helm charts
gh workflow run ci-helm.yaml

# Create release (dry run mode)
gh workflow run release.yaml -f dry_run=true

# Trigger Renovate dependency check
gh workflow run renovate.yaml
```

### Workflow Status

```bash
# List recent workflow runs
gh run list

# View specific run details
gh run view <run-id>

# Watch running workflow
gh run watch <run-id>

# Download artifacts
gh run download <run-id>
```

### Image Management

```bash
# List images in GHCR
gh api user/packages/container/test-workflow%2Fapi-gateway/versions

# Delete old image versions (keep last 10)
gh api -X DELETE user/packages/container/test-workflow%2Fapi-gateway/versions/<version-id>
```

---

## References

### GitHub Documentation

- GitHub Actions: https://docs.github.com/en/actions
- GitHub Container Registry: https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry
- GitHub Security Features: https://docs.github.com/en/code-security

### CI/CD Tools

- Semantic Release: https://semantic-release.gitbook.io/semantic-release/
- Conventional Commits: https://www.conventionalcommits.org/
- Renovate: https://docs.renovatebot.com/

### Security

- Trivy Documentation: https://aquasecurity.github.io/trivy/
- SARIF Format: https://sarifweb.azurewebsites.net/

### Container Tools

- Docker BuildKit: https://docs.docker.com/build/buildkit/
- Multi-stage Builds: https://docs.docker.com/build/building/multi-stage/

### Helm

- Helm Documentation: https://helm.sh/docs/
- Chart Best Practices: https://helm.sh/docs/chart_best_practices/

### Related Project Documentation

- [Renovate Configuration](./docs/RENOVATE.md)
- [Semantic Release Guide](./docs/SEMANTIC_RELEASE.md)
- [Trivy Security Scanning](./docs/TRIVY.md)
