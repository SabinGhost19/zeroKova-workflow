# 🔧 GitHub Actions CI/CD Pipeline Guide

This document explains how to configure and use the CI/CD pipeline for the test-workflow project.

## 📁 Pipeline Structure

The pipeline is split into separate workflow files for better maintainability:

```
.github/workflows/
├── ci-test.yaml         # Testing for all microservices
├── ci-build.yaml        # Build and push Docker images
├── ci-security.yaml     # Trivy security scanning
├── ci-helm.yaml         # Helm chart validation
├── release.yaml         # Semantic release automation
├── renovate.yaml        # Dependency updates
└── pipeline-summary.yaml # Aggregated status
```

## ⚙️ Required Configuration

### 1. Repository Settings

Navigate to **Settings → Actions → General**:

- [x] Allow all actions and reusable workflows
- [x] Read and write permissions (under "Workflow permissions")
- [x] Allow GitHub Actions to create and approve pull requests

### 2. Required Secrets

No additional secrets are required! The pipeline uses `GITHUB_TOKEN` which is automatically provided.

| Secret | Source | Description |
|--------|--------|-------------|
| `GITHUB_TOKEN` | Automatic | Used for GHCR push, releases, and API calls |

### 3. Package Registry (GHCR)

Images are pushed to GitHub Container Registry (GHCR). Ensure:

1. Go to **Settings → Packages**
2. Verify "Inherit access from source repository" is enabled

### 4. Branch Protection (Recommended)

For the `main` branch:

1. Go to **Settings → Branches → Add rule**
2. Configure:
   - [x] Require a pull request before merging
   - [x] Require status checks to pass
   - [x] Select: `Test api-gateway`, `Test order-service`, etc.
   - [x] Require branches to be up to date
   - [x] Include administrators

## 🔄 Workflow Triggers

| Workflow | Push to main | Pull Request | Manual | Schedule | Workflow Run |
|----------|:------------:|:------------:|:------:|:--------:|:------------:|
| CI - Tests | ✅ | ✅ | ✅ | - | - |
| CI - Build & Push | ✅ | - | ✅ | - | - |
| CI - Security Scan | - | - | ✅ | - | ✅ (after build) |
| CI - Helm Validation | ✅ | ✅ | ✅ | - | - |
| Release | ✅ | - | ✅ | - | - |
| Renovate | - | - | ✅ | Daily 6AM UTC | - |
| Pipeline Summary | - | - | ✅ | - | ✅ |

## 📦 Docker Images

### Image Naming

All images are published to:
```
ghcr.io/{owner}/test-workflow/{service}:{tag}
```

### Available Tags

| Tag Pattern | Description |
|-------------|-------------|
| `latest` | Latest build from main branch |
| `sha-{commit}` | Specific commit build |
| `v{x.y.z}` | Semantic version (created by release) |
| `{branch}` | Branch name (for non-main branches) |

### Pulling Images

```bash
# Pull latest
docker pull ghcr.io/sabinghosty19/test-workflow/api-gateway:latest

# Pull specific version
docker pull ghcr.io/sabinghosty19/test-workflow/api-gateway:v1.0.0
```

## 🔒 Security Scanning

### Trivy Reports Location

After each scan, reports are available in:

1. **GitHub Actions → Artifacts**
   - `trivy-security-reports-all` - Combined report with index.html
   - Individual reports per service

2. **GitHub Security Tab**
   - **Security → Code scanning alerts**
   - View all vulnerabilities in one place

### Report Types

| Format | Purpose |
|--------|---------|
| HTML | Human-readable, downloadable |
| SARIF | GitHub Security integration |
| JSON | Machine-readable processing |

## 🚀 Releasing

### Automatic Releases

Releases are triggered automatically when:
1. Push to `main` branch
2. Commits follow [Conventional Commits](https://www.conventionalcommits.org/)

### Version Bumping

| Commit Type | Version Bump | Example |
|-------------|--------------|---------|
| `fix:` | Patch (0.0.X) | `fix: resolve null pointer` |
| `feat:` | Minor (0.X.0) | `feat: add new endpoint` |
| `feat!:` or `BREAKING CHANGE:` | Major (X.0.0) | `feat!: change API format` |

### What Happens on Release

1. Version calculated from commits
2. `CHANGELOG.md` updated
3. Git tag created (e.g., `v1.2.3`)
4. GitHub Release created
5. Docker images tagged with version

## 🔄 Dependency Updates (Renovate)

### Automatic Behavior

- Runs daily at 6:00 AM UTC
- Creates PRs for dependency updates
- **Automerges** patch and minor updates
- **Requires review** for major updates

### PR Labels

| Label | Meaning |
|-------|---------|
| `dependencies` | All dependency PRs |
| `renovate` | Created by Renovate |
| `major-update` | Major version bump |
| `security` | Security-related update |

### Skipping Updates

To skip a specific update, add to the PR comment:
```
@renovate ignore this dependency
```

## 🐛 Troubleshooting

### Tests Failing

1. Check the specific test job in Actions
2. Look for error messages in logs
3. Run tests locally:
   ```bash
   # Go
   cd services/api-gateway && go test ./...
   
   # Java
   cd services/order-service && mvn test
   
   # C#
   cd services/inventory-service && dotnet test
   
   # Ruby
   cd services/notification-service && bundle exec rspec
   
   # Node
   cd frontend && npm test
   ```

### Build Failing

1. Verify Dockerfile syntax
2. Check build logs for errors
3. Test locally:
   ```bash
   docker build -t test services/api-gateway
   ```

### Security Scan Issues

1. Scans continue even with vulnerabilities
2. Check HTML reports in artifacts
3. Review Security tab for details

### Release Not Triggering

1. Ensure commits follow Conventional Commits format
2. Check for `[skip ci]` in commit messages
3. Verify branch is `main`

## 📋 Quick Commands

```bash
# Manual trigger - Tests
gh workflow run ci-test.yaml

# Manual trigger - Build all
gh workflow run ci-build.yaml -f force_build_all=true

# Manual trigger - Security scan
gh workflow run ci-security.yaml

# Manual trigger - Release (dry run)
gh workflow run release.yaml -f dry_run=true

# Manual trigger - Renovate
gh workflow run renovate.yaml
```

## 📚 Related Documentation

- [Renovate Configuration](./docs/RENOVATE.md)
- [Semantic Release Guide](./docs/SEMANTIC_RELEASE.md)
- [Trivy Security Scanning](./docs/TRIVY.md)
