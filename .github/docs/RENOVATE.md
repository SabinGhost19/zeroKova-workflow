# 🔄 Renovate - Automated Dependency Updates

Renovate is a bot that keeps your dependencies up-to-date by automatically creating pull requests when new versions are available.

## 📋 Table of Contents

- [What It Does](#what-it-does)
- [Configuration](#configuration)
- [Behavior](#behavior)
- [Package Rules](#package-rules)
- [Automerge](#automerge)
- [Grouping](#grouping)
- [Commands](#commands)
- [Troubleshooting](#troubleshooting)

## 🎯 What It Does

Renovate monitors your project's dependencies and:

1. **Detects outdated dependencies** across all supported package managers
2. **Creates Pull Requests** with updates
3. **Automerges** non-breaking updates (patch/minor)
4. **Groups updates** by language/type to reduce PR noise
5. **Handles security updates** with priority

### Supported Package Managers

| Manager | Files | Services |
|---------|-------|----------|
| Go modules | `go.mod`, `go.sum` | api-gateway |
| Maven | `pom.xml` | order-service |
| NuGet | `*.csproj` | inventory-service |
| Bundler | `Gemfile`, `Gemfile.lock` | notification-service |
| npm | `package.json`, `package-lock.json` | frontend |
| Docker | `Dockerfile` | All services |
| Helm | `Chart.yaml`, `values.yaml` | helm/test-workflow |
| GitHub Actions | `.github/workflows/*.yaml` | CI/CD |

## ⚙️ Configuration

The configuration is in `renovate.json` at the project root.

### Key Settings

```json
{
  "schedule": ["every weekday"],    // When to run
  "timezone": "Europe/Bucharest",   // Your timezone
  "automerge": true,                // Automerge enabled
  "prHourlyLimit": 5,               // Max PRs per hour
  "prConcurrentLimit": 10,          // Max open PRs
  "minimumReleaseAge": "3 days"     // Wait before updating
}
```

## 📅 Behavior

### Schedule

| Action | Schedule |
|--------|----------|
| Check for updates | Daily at 6:00 AM UTC |
| Create PRs | Weekdays only |
| Lock file maintenance | Monday before 6 AM |

### Update Flow

```
New version published
        ↓
Wait 3 days (minimumReleaseAge)
        ↓
Renovate creates PR
        ↓
    ┌───────────┐
    │ Patch/Minor │ → Automerge ✅
    └───────────┘
    ┌───────────┐
    │   Major   │ → Manual review required ❌
    └───────────┘
```

## 📦 Package Rules

### By Update Type

| Type | Automerge | Labels |
|------|-----------|--------|
| Patch (0.0.X) | ✅ Yes | `dependencies`, `renovate` |
| Minor (0.X.0) | ✅ Yes | `dependencies`, `renovate` |
| Major (X.0.0) | ❌ No | `dependencies`, `major-update`, `needs-review` |
| Security | ✅ Yes (immediate) | `security`, `dependencies` |

### By Language

All non-major updates are grouped by language:

| Group | Included |
|-------|----------|
| `go-dependencies` | All Go module updates |
| `npm-dependencies` | All npm package updates |
| `maven-dependencies` | All Maven artifact updates |
| `nuget-dependencies` | All NuGet package updates |
| `ruby-dependencies` | All Ruby gem updates |
| `docker-images` | All Dockerfile base image updates |
| `github-actions` | All GitHub Actions updates |
| `helm-charts` | All Helm chart updates |

## 🔀 Automerge

### When Automerge Happens

1. ✅ All CI checks pass
2. ✅ Update is patch or minor
3. ✅ No merge conflicts
4. ✅ Not a major version bump

### Automerge Strategy

```json
{
  "automergeType": "pr",
  "platformAutomerge": true,
  "automergeStrategy": "squash"
}
```

- Creates a PR (for visibility)
- Uses GitHub's native automerge
- Squashes commits for clean history

## 📊 Grouping

Instead of one PR per dependency, updates are grouped:

### Before Grouping
```
PR #1: Update axios to 1.5.0
PR #2: Update react to 18.2.1
PR #3: Update typescript to 5.2.0
...
```

### After Grouping
```
PR #1: chore(deps): update npm-dependencies
       - axios 1.4.0 → 1.5.0
       - react 18.2.0 → 18.2.1
       - typescript 5.1.0 → 5.2.0
```

## 🎮 Commands

You can interact with Renovate using comments on PRs:

### Dependency Dashboard

A special issue called "Dependency Dashboard" provides an overview. You can:
- Check/uncheck dependencies to update
- See pending updates
- Force immediate updates

### PR Commands

Comment on a Renovate PR:

| Command | Effect |
|---------|--------|
| `@renovate rebase` | Rebase the PR |
| `@renovate recreate` | Close and recreate the PR |
| `@renovate refresh` | Update the PR |
| `@renovate ignore this dependency` | Never update this dependency |
| `@renovate ignore this major version` | Skip this major version |
| `@renovate unignore` | Remove ignore rules |

## 🔒 Security Updates

Security updates receive special treatment:

```json
{
  "matchCategories": ["security"],
  "automerge": true,
  "minimumReleaseAge": "0 days",
  "prPriority": 10
}
```

- **No waiting period** - immediate PR
- **Highest priority** - processed first
- **Automerge enabled** - if tests pass
- **Special label** - `security`

## 🐛 Troubleshooting

### PRs Not Being Created

1. Check if Renovate workflow ran:
   ```
   Actions → Renovate → Recent runs
   ```

2. Check Dependency Dashboard issue for errors

3. Verify configuration:
   ```bash
   npx renovate-config-validator
   ```

### Automerge Not Working

1. Ensure branch protection allows automerge
2. Check all required status checks pass
3. Verify no merge conflicts exist

### Too Many PRs

Adjust limits in `renovate.json`:
```json
{
  "prHourlyLimit": 2,
  "prConcurrentLimit": 5
}
```

### Ignoring a Package

Add to `renovate.json`:
```json
{
  "ignoreDeps": ["problematic-package"]
}
```

Or comment on PR:
```
@renovate ignore this dependency
```

## 📚 Resources

- [Renovate Documentation](https://docs.renovatebot.com/)
- [Configuration Options](https://docs.renovatebot.com/configuration-options/)
- [Presets](https://docs.renovatebot.com/presets-config/)
- [Regex Managers](https://docs.renovatebot.com/modules/manager/regex/)
