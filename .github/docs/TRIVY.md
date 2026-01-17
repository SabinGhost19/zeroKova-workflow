# 🔍 Trivy - Security Scanning

Trivy is a comprehensive security scanner that detects vulnerabilities in container images, filesystems, and Infrastructure as Code (IaC).

## 📋 Table of Contents

- [What It Does](#what-it-does)
- [Scan Types](#scan-types)
- [Reports](#reports)
- [Configuration](#configuration)
- [Viewing Results](#viewing-results)
- [Understanding Vulnerabilities](#understanding-vulnerabilities)
- [Best Practices](#best-practices)

## 🎯 What It Does

Trivy scans for:

1. **Container Image Vulnerabilities** - CVEs in OS packages and application dependencies
2. **IaC Misconfigurations** - Security issues in Helm, Kubernetes, and Flux configs
3. **Secret Detection** - Accidentally committed secrets and credentials

### What Gets Scanned

| Target | Content |
|--------|---------|
| Docker Images | All 5 microservice images |
| Helm Charts | `helm/test-workflow/` |
| K8s Manifests | `k8s/` directory |
| Flux Configs | `flux/` directory |

## 🔬 Scan Types

### 1. Container Image Scanning

Scans built Docker images for:
- OS package vulnerabilities (apt, yum, apk)
- Language-specific vulnerabilities (npm, pip, go, maven, nuget, bundler)
- Base image vulnerabilities

```yaml
- image: ghcr.io/sabinghosty19/test-workflow/api-gateway:latest
  scanned-for:
    - CVEs in Alpine/Debian base image
    - CVEs in Go dependencies
```

### 2. IaC (Infrastructure as Code) Scanning

Scans configuration files for:
- Missing security contexts
- Privileged containers
- Missing resource limits
- Insecure network policies
- Secrets in plaintext

```yaml
scanned-paths:
  - helm/
  - k8s/
  - flux/
```

## 📊 Reports

### Report Formats

| Format | Purpose | Location |
|--------|---------|----------|
| **HTML** | Human-readable, downloadable | GitHub Artifacts |
| **SARIF** | GitHub Security integration | Security tab → Code scanning |
| **JSON** | Machine processing | GitHub Artifacts |

### Accessing Reports

#### Method 1: GitHub Artifacts

1. Go to **Actions** tab
2. Click on a workflow run
3. Scroll to **Artifacts** section
4. Download `trivy-security-reports-all`
5. Open `index.html` for overview

#### Method 2: GitHub Security Tab

1. Go to **Security** tab
2. Click **Code scanning**
3. View all alerts grouped by severity

### Report Structure

```
trivy-security-reports-all/
├── index.html                          # Overview page
├── trivy-report-api-gateway/
│   ├── trivy-api-gateway-vuln.html     # Image vulnerabilities
│   └── trivy-api-gateway-vuln.json
├── trivy-report-order-service/
│   └── ...
├── trivy-report-inventory-service/
│   └── ...
├── trivy-report-notification-service/
│   └── ...
├── trivy-report-frontend/
│   └── ...
└── trivy-iac-reports/
    ├── trivy-helm-misconfig.html       # Helm misconfigurations
    ├── trivy-k8s-misconfig.html        # K8s misconfigurations
    └── trivy-flux-misconfig.html       # Flux misconfigurations
```

## ⚙️ Configuration

### Current Settings

| Setting | Value | Description |
|---------|-------|-------------|
| Severity | CRITICAL, HIGH, MEDIUM, LOW | All severities scanned |
| Ignore Unfixed | ✅ Yes | Skip vulnerabilities without patches |
| Exit Code | 0 | Pipeline continues on findings |
| Retention | 90 days | How long artifacts are kept |

### Why Ignore Unfixed?

Unfixed vulnerabilities have no available patch. Reporting them creates noise without actionable fixes. Once a fix is released, they will appear in scans.

## 👀 Viewing Results

### Summary in Workflow

Each scan provides a summary in the Actions tab:

```markdown
### 🔍 Trivy Scan: api-gateway

| Severity | Count |
|----------|-------|
| 🔴 Critical | 0 |
| 🟠 High | 2 |
| 🟡 Medium | 5 |
| 🟢 Low | 12 |

> ℹ️ Unfixed vulnerabilities are ignored
```

### HTML Report

The HTML report provides:
- Vulnerability details
- CVE identifiers
- Affected packages
- Fixed versions
- CVSS scores
- Links to advisories

### GitHub Security Tab

Benefits:
- Centralized view of all vulnerabilities
- Filter by severity, package, or file
- Track resolution status
- Get notifications on new findings

## 🎯 Understanding Vulnerabilities

### Severity Levels

| Level | CVSS Score | Description | Action |
|-------|------------|-------------|--------|
| 🔴 Critical | 9.0 - 10.0 | Immediate risk | Fix ASAP |
| 🟠 High | 7.0 - 8.9 | Significant risk | Fix soon |
| 🟡 Medium | 4.0 - 6.9 | Moderate risk | Plan to fix |
| 🟢 Low | 0.1 - 3.9 | Minor risk | Low priority |

### Vulnerability Details

Each finding includes:

```
┌─────────────────────────────────────────────────────────────┐
│ CVE-2024-12345                                              │
├─────────────────────────────────────────────────────────────┤
│ Package:    openssl                                         │
│ Version:    1.1.1k                                          │
│ Fixed In:   1.1.1l                                          │
│ Severity:   HIGH (7.5)                                      │
│ Title:      Buffer overflow in X509 certificate parsing     │
│ References: https://nvd.nist.gov/vuln/detail/CVE-2024-12345 │
└─────────────────────────────────────────────────────────────┘
```

## 🛠️ Fixing Vulnerabilities

### Container Images

1. **Update base image**
   ```dockerfile
   # Before
   FROM node:18-alpine3.17
   
   # After (use latest patch)
   FROM node:18-alpine3.19
   ```

2. **Update dependencies**
   ```bash
   # npm
   npm update
   npm audit fix
   
   # Go
   go get -u ./...
   
   # Maven
   mvn versions:use-latest-releases
   ```

3. **Rebuild and push**
   - Push to main triggers rebuild
   - New scan will verify fixes

### IaC Misconfigurations

Common fixes:

```yaml
# Add security context
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL

# Add resource limits
resources:
  limits:
    cpu: "500m"
    memory: "256Mi"
  requests:
    cpu: "100m"
    memory: "128Mi"

# Add network policy
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
```

## ✅ Best Practices

### 1. Use Minimal Base Images

```dockerfile
# ❌ Bad - full OS with many packages
FROM ubuntu:22.04

# ✅ Good - minimal image
FROM alpine:3.19

# ✅ Better - distroless (no shell, minimal attack surface)
FROM gcr.io/distroless/static
```

### 2. Multi-stage Builds

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app

# Runtime stage - minimal image
FROM alpine:3.19
COPY --from=builder /app/app /app
USER 1000
ENTRYPOINT ["/app"]
```

### 3. Pin Versions

```dockerfile
# ❌ Bad - unpredictable versions
FROM node:latest

# ✅ Good - specific version
FROM node:20.11.0-alpine3.19
```

### 4. Regular Updates

- Keep base images updated
- Update dependencies regularly (Renovate helps!)
- Schedule periodic rebuilds

### 5. Review Before Ignoring

If you must ignore a vulnerability:

1. Document why it's acceptable
2. Set a reminder to revisit
3. Consider compensating controls

## 🐛 Troubleshooting

### No Scan Results

1. Verify image exists in GHCR
2. Check if build workflow succeeded
3. Review workflow logs for errors

### Too Many Vulnerabilities

1. Update base image
2. Remove unused dependencies
3. Consider switching to minimal/distroless images

### False Positives

Some vulnerabilities may not apply to your usage. To ignore:

1. Create `.trivyignore` file:
   ```
   # Ignore specific CVE
   CVE-2024-12345
   
   # Ignore with comment
   CVE-2024-67890  # Not applicable - feature not used
   ```

2. Add to repository root

## 📚 Resources

- [Trivy Documentation](https://aquasecurity.github.io/trivy/)
- [Trivy GitHub](https://github.com/aquasecurity/trivy)
- [CVE Database](https://cve.mitre.org/)
- [NIST NVD](https://nvd.nist.gov/)
