# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.10](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.9...v1.3.10) (2026-01-18)

### 🐛 Bug Fixes

* **proto:** enable java_multiple_files for correct class generation ([ea528e2](https://github.com/SabinGhost19/zeroKova-workflow/commit/ea528e26ff7ab440d84fcf31368b83b85b7150a0))

## [1.3.9](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.8...v1.3.9) (2026-01-18)

### 🐛 Bug Fixes

* **order-service:** ensure protobuf compilation runs in generate-sources ([e1b8896](https://github.com/SabinGhost19/zeroKova-workflow/commit/e1b889669cfcab3e631ede9ad2a592684fe726f0))

## [1.3.8](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.7...v1.3.8) (2026-01-18)

### 🐛 Bug Fixes

* downgrade typescript to 4.9.5 to match react-scripts requirements ([1be8971](https://github.com/SabinGhost19/zeroKova-workflow/commit/1be8971519d8c3b4d0003d6cfa79ef4c5b863931))

## [1.3.7](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.6...v1.3.7) (2026-01-18)

### 🐛 Bug Fixes

* add protobuf source generation configuration for order-service ([54e8bdd](https://github.com/SabinGhost19/zeroKova-workflow/commit/54e8bddc254248b78435bdfe5b2d5a673d52ff00))

## [1.3.6](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.5...v1.3.6) (2026-01-18)

### 🐛 Bug Fixes

* disambiguate Product reference in inventory service ([3947398](https://github.com/SabinGhost19/zeroKova-workflow/commit/39473982824c2717335f4e6b85c5e066fcc23205))

## [1.3.5](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.4...v1.3.5) (2026-01-18)

### 🐛 Bug Fixes

* fix C# inventory service - correct proto namespace and add xunit ([feedd8a](https://github.com/SabinGhost19/zeroKova-workflow/commit/feedd8a4602435d687a63183017e2df55df3c60f))
* fix C# inventory service - correct proto namespace and add xunit ([3d4d49c](https://github.com/SabinGhost19/zeroKova-workflow/commit/3d4d49c811a86d1c097288037870ed6c91537b75))

## [1.3.4](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.3...v1.3.4) (2026-01-18)

### 🐛 Bug Fixes

* generate all lock files in CI, remove local dependencies ([7b41f65](https://github.com/SabinGhost19/zeroKova-workflow/commit/7b41f65c08152c3403fd396da545a1d95d572a24))
* generate lock files in CI instead of committing them locally ([7061f2f](https://github.com/SabinGhost19/zeroKova-workflow/commit/7061f2fa7476734b9f20a55f8d88d576f341ac49))

## [1.3.3](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.2...v1.3.3) (2026-01-18)

### 🐛 Bug Fixes

* CI pipeline fixes and add .gitignore ([#6](https://github.com/SabinGhost19/zeroKova-workflow/issues/6)) ([aafa29f](https://github.com/SabinGhost19/zeroKova-workflow/commit/aafa29f33eb5d2093b1855f5964f4225eda9cd0b))

## [1.3.2](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.1...v1.3.2) (2026-01-18)

### 🐛 Bug Fixes

* CI pipeline - add Linux platforms to Gemfile.lock, fix ESLint warning, disable cancel-in-progress ([6d010d8](https://github.com/SabinGhost19/zeroKova-workflow/commit/6d010d8f02ab1f65fb7882c8e13ee535cfb51130))

## [1.3.1](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.3.0...v1.3.1) (2026-01-18)

### 🐛 Bug Fixes

* simplify CI pipeline and resolve build errors ([10d0f08](https://github.com/SabinGhost19/zeroKova-workflow/commit/10d0f08c96dfca849e8b6981192aa9131bef6f19))

## [1.3.0](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.2.0...v1.3.0) (2026-01-18)

### 🚀 Features

* add version comments to all microservices ([0d0071f](https://github.com/SabinGhost19/zeroKova-workflow/commit/0d0071ff03aed9c3dbe83defe311d7ec1c8f64d2))

## [1.2.0](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.1.1...v1.2.0) (2026-01-18)

### 🚀 Features

* improve manual trigger for CI workflows ([7055c7f](https://github.com/SabinGhost19/zeroKova-workflow/commit/7055c7f6575b070a45bbcb40d55329bdf849c964))

## [1.1.1](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.1.0...v1.1.1) (2026-01-18)

### 🐛 Bug Fixes

* regenerate go.sum to fix checksum mismatch ([6febd81](https://github.com/SabinGhost19/zeroKova-workflow/commit/6febd81bd701e78439dc6332213b2f8abf53e087))

## [1.1.0](https://github.com/SabinGhost19/zeroKova-workflow/compare/v1.0.0...v1.1.0) (2026-01-18)

### 🚀 Features

* add api-gateway header comment ([2ac08a2](https://github.com/SabinGhost19/zeroKova-workflow/commit/2ac08a264f8837652698286536a08fb67dbc5087))

## 1.0.0 (2026-01-18)

### 🐛 Bug Fixes

* ghrc image prefix ([1cb722b](https://github.com/SabinGhost19/zeroKova-workflow/commit/1cb722b9053cd52be54ff06354fe086920c1f818))
