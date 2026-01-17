# 🚀 Semantic Release - Automated Versioning & Releases

Semantic Release automates the entire package release workflow including determining the next version number, generating release notes, and publishing to GitHub Releases.

## 📋 Table of Contents

- [What It Does](#what-it-does)
- [How It Works](#how-it-works)
- [Conventional Commits](#conventional-commits)
- [Version Bumping](#version-bumping)
- [Changelog](#changelog)
- [Git Tags](#git-tags)
- [Best Practices](#best-practices)
- [Examples](#examples)

## 🎯 What It Does

On every push to `main`, Semantic Release:

1. **Analyzes commits** since the last release
2. **Determines version bump** (major/minor/patch)
3. **Updates CHANGELOG.md** with release notes
4. **Creates git tag** (e.g., `v1.2.3`)
5. **Creates GitHub Release** with notes
6. **Tags Docker images** with version

## ⚙️ How It Works

```
Push to main
     ↓
Analyze commits since last tag
     ↓
┌─────────────────────────────┐
│ Found release-worthy commits?│
└─────────────────────────────┘
     │                    │
    Yes                   No
     ↓                    ↓
Calculate next version   Exit
     ↓
Update CHANGELOG.md
     ↓
Create git tag (v1.2.3)
     ↓
Create GitHub Release
     ↓
Tag Docker images
```

## 📝 Conventional Commits

### Commit Message Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New feature | **Minor** (0.X.0) |
| `fix` | Bug fix | **Patch** (0.0.X) |
| `perf` | Performance improvement | Patch |
| `revert` | Revert previous commit | Patch |
| `refactor` | Code refactoring | Patch |
| `docs` | Documentation only | ❌ No release |
| `style` | Code style (formatting) | ❌ No release |
| `test` | Adding tests | ❌ No release |
| `build` | Build system changes | ❌ No release |
| `ci` | CI configuration | ❌ No release |
| `chore` | Maintenance | ❌ No release |

### Breaking Changes

For **Major** version bumps (X.0.0), use:

```bash
# Option 1: Add ! after type
feat!: change authentication API

# Option 2: Footer with BREAKING CHANGE
feat: change authentication API

BREAKING CHANGE: The authentication endpoint now requires OAuth2
```

### Scope (Optional)

Scope indicates what part of the code is affected:

```bash
feat(api-gateway): add rate limiting
fix(order-service): resolve null pointer
docs(readme): update installation steps
```

## 📊 Version Bumping

### Rules

| Commits Since Last Release | Next Version |
|---------------------------|--------------|
| Only `docs:`, `chore:`, `ci:`, `test:` | No release |
| At least one `fix:` | Patch (0.0.X) |
| At least one `feat:` | Minor (0.X.0) |
| Any `BREAKING CHANGE` or `!` | Major (X.0.0) |

### Examples

Current version: `v1.2.3`

| Commits | Next Version |
|---------|--------------|
| `fix: resolve bug`, `docs: update readme` | v1.2.4 |
| `feat: add new endpoint` | v1.3.0 |
| `feat!: change API format` | v2.0.0 |
| `fix: bug`, `feat: feature`, `docs: readme` | v1.3.0 |

## 📰 Changelog

A `CHANGELOG.md` is automatically generated and updated with each release.

### Format

```markdown
# Changelog

## [1.3.0](https://github.com/user/repo/compare/v1.2.3...v1.3.0) (2024-01-15)

### 🚀 Features

* **api-gateway:** add rate limiting ([abc1234](link))
* **order-service:** implement order tracking ([def5678](link))

### 🐛 Bug Fixes

* **frontend:** fix responsive layout ([ghi9012](link))

## [1.2.3](https://github.com/user/repo/compare/v1.2.2...v1.2.3) (2024-01-10)

### 🐛 Bug Fixes

* **inventory-service:** resolve stock calculation ([jkl3456](link))
```

### Sections

| Section | Included Types |
|---------|---------------|
| 🚀 Features | `feat` |
| 🐛 Bug Fixes | `fix` |
| ⚡ Performance | `perf` |
| ⏪ Reverts | `revert` |
| ♻️ Code Refactoring | `refactor` |
| 📚 Documentation | `docs` |

## 🏷️ Git Tags

### Tag Format

Tags follow semantic versioning: `v{major}.{minor}.{patch}`

```
v1.0.0
v1.0.1
v1.1.0
v2.0.0
```

### Docker Image Tags

After a release, Docker images are tagged with the version:

```
ghcr.io/sabinghosty19/test-workflow/api-gateway:v1.3.0
ghcr.io/sabinghosty19/test-workflow/api-gateway:latest
```

## ✅ Best Practices

### Do's ✅

```bash
# Good commit messages
feat(api): add user authentication endpoint
fix(order): resolve null pointer when order is empty
docs: update API documentation
refactor(inventory): simplify stock calculation logic
perf(frontend): optimize image loading

# Breaking change with explanation
feat(api)!: change response format to JSON:API

BREAKING CHANGE: All API responses now follow JSON:API specification.
Migration guide: https://example.com/migration
```

### Don'ts ❌

```bash
# Bad commit messages
update code                    # Not descriptive
fixed bug                      # Missing type
feat: added stuff              # Vague description
WIP                            # Work in progress shouldn't be on main
fix bug in the thing          # Missing type and colon
```

### Atomic Commits

Each commit should:
- Focus on one logical change
- Be complete and working
- Have a clear, descriptive message

```bash
# Instead of:
git commit -m "fix: various fixes and updates"

# Do:
git commit -m "fix(auth): resolve token expiration issue"
git commit -m "feat(api): add pagination to list endpoints"
git commit -m "docs: update authentication guide"
```

## 📖 Examples

### Feature Development

```bash
# Adding a new feature
git commit -m "feat(order-service): add order cancellation"

# Result: v1.2.0 → v1.3.0
```

### Bug Fix

```bash
# Fixing a bug
git commit -m "fix(inventory): correct stock level calculation"

# Result: v1.3.0 → v1.3.1
```

### Breaking Change

```bash
# Changing API format
git commit -m "feat(api)!: migrate to GraphQL

BREAKING CHANGE: REST API endpoints are deprecated.
Use GraphQL endpoint at /graphql instead."

# Result: v1.3.1 → v2.0.0
```

### Multiple Commits

```bash
# Multiple commits in one PR
git commit -m "feat(frontend): add dark mode toggle"
git commit -m "fix(frontend): fix mobile navigation"
git commit -m "docs: update theming guide"

# Result: v1.3.0 → v1.4.0 (feat takes precedence)
```

### No Release

```bash
# Documentation only
git commit -m "docs: improve README installation steps"

# Result: No new release
```

## 🛠️ Configuration

The configuration is in `.releaserc.json`:

```json
{
  "branches": ["main"],
  "tagFormat": "v${version}",
  "plugins": [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/changelog",
    "@semantic-release/github",
    "@semantic-release/git"
  ]
}
```

## 🔍 Dry Run

Test what would happen without actually releasing:

```bash
# Via GitHub Actions
gh workflow run release.yaml -f dry_run=true

# Locally
npx semantic-release --dry-run
```

## 🐛 Troubleshooting

### No Release Created

1. Check commit messages follow Conventional Commits
2. Ensure commits include release-worthy types (`feat:`, `fix:`)
3. Verify pushing to `main` branch
4. Check workflow logs for errors

### Wrong Version Bump

1. Review commit messages for unexpected types
2. Check for unintended `BREAKING CHANGE` footers
3. Verify commit message format is correct

### Changelog Not Updated

1. Ensure `@semantic-release/changelog` plugin is configured
2. Check `.releaserc.json` for correct plugin order
3. Verify git permissions for committing changes

## 📚 Resources

- [Semantic Release](https://semantic-release.gitbook.io/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
