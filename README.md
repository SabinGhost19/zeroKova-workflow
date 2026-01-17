# Test Workflow - Zero Trust Microservices Demo

A polyglot microservices application demonstrating zero trust security patterns in Kubernetes. This project implements a complete e-commerce order management system with multiple backend services, each written in a different programming language.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              FRONTEND (React + Nginx)                    │
│                                   Port: 80                               │
└─────────────────────────────────────┬───────────────────────────────────┘
                                      │ REST
                                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           API GATEWAY (Go + Gin)                         │
│                              Port: 8080                                  │
└────────┬────────────────────────────┼────────────────────────┬──────────┘
         │ gRPC                       │ gRPC                   │ gRPC
         ▼                            ▼                        ▼
┌─────────────────┐      ┌─────────────────────┐    ┌─────────────────────┐
│  ORDER SERVICE  │      │  INVENTORY SERVICE  │    │ NOTIFICATION SERVICE│
│ (Java/Spring)   │──────│    (C# .NET)        │    │      (Ruby)         │
│  Port: 50051    │Kafka │    Port: 50052      │    │    Port: 50053      │
└────────┬────────┘      └──────────┬──────────┘    └──────────┬──────────┘
         │                          │                          │
         └──────────────────────────┼──────────────────────────┘
                                    │
                                    ▼
                    ┌───────────────────────────────┐
                    │        PostgreSQL             │
                    │        Port: 5432             │
                    │   (Schema per service)        │
                    └───────────────────────────────┘
```

## Features

- **Multi-language microservices**: Go, Java, C#, Ruby, React
- **gRPC communication**: Efficient inter-service communication using protocol buffers
- **Event-driven architecture**: Kafka for asynchronous messaging between services
- **Zero Trust security**: Network policies enforcing least-privilege access
- **GitOps deployment**: Flux CD integration for continuous deployment
- **Helm packaging**: Single chart for all microservices with modular configuration
- **CI/CD pipeline**: GitHub Actions with matrix builds and GHCR integration

## Components

| Service | Language | HTTP Port | gRPC Port | Description |
|---------|----------|-----------|-----------|-------------|
| Frontend | React/TypeScript | 80 | - | User interface served by Nginx |
| API Gateway | Go/Gin | 8080 | - | REST entry point, routes to backend services |
| Order Service | Java/Spring Boot | 8081 | 50051 | Order management, Kafka producer |
| Inventory Service | C# .NET 8 | 8082 | 50052 | Stock management, Kafka consumer |
| Notification Service | Ruby | 8083 | 50053 | Notification handling |

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Kubernetes cluster (minikube, kind, or cloud provider)
- Helm 3.x
- kubectl

### Local Development (Docker Compose)

```bash
# build all images
./scripts/build-all.sh

# start the application
docker-compose up -d

# verify services are running
docker-compose ps

# access the frontend
open http://localhost:3000
```

### Kubernetes Deployment

#### Using Helm

```bash
# add helm dependencies
cd helm/test-workflow
helm dependency update

# install the chart
helm install test-workflow ./helm/test-workflow \
  --namespace test-workflow \
  --create-namespace

# verify deployment
kubectl get pods -n test-workflow
```

#### Using Kustomize

```bash
# deploy using the provided script
./scripts/deploy-k8s.sh

# or manually with kustomize
kubectl apply -k k8s/base

# add to /etc/hosts for local access
echo "127.0.0.1 test-workflow.local" | sudo tee -a /etc/hosts
```

### GitOps with Flux CD

```bash
# bootstrap flux (if not already installed)
flux bootstrap github \
  --owner=sabinghosty19 \
  --repository=test-workflow \
  --path=flux/overlays/development \
  --personal

# check deployment status
flux get helmreleases -n flux-system
```

## Project Structure

```
test-workflow/
├── .github/
│   ├── README.md                 # pipeline configuration guide
│   ├── docs/
│   │   ├── RENOVATE.md           # renovate documentation
│   │   ├── SEMANTIC_RELEASE.md   # semantic release guide
│   │   └── TRIVY.md              # trivy security scanning guide
│   └── workflows/
│       ├── ci-test.yaml          # testing workflow
│       ├── ci-build.yaml         # build and push workflow
│       ├── ci-security.yaml      # trivy security scanning
│       ├── ci-helm.yaml          # helm validation
│       ├── release.yaml          # semantic release
│       ├── renovate.yaml         # dependency updates
│       └── pipeline-summary.yaml # status aggregation
├── proto/                        # grpc protobuf definitions
│   ├── order.proto
│   ├── inventory.proto
│   ├── notification.proto
│   └── common.proto
├── services/
│   ├── api-gateway/              # go api gateway
│   ├── order-service/            # java spring boot service
│   ├── inventory-service/        # c# .net service
│   └── notification-service/     # ruby service
├── frontend/                     # react spa
├── helm/
│   └── test-workflow/            # helm chart
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
├── flux/                         # flux cd configuration
│   ├── base/
│   └── overlays/
│       ├── development/
│       ├── staging/
│       └── production/
├── k8s/                          # kubernetes manifests
│   ├── base/
│   └── services/
├── scripts/                      # helper scripts
├── renovate.json                 # renovate configuration
├── .releaserc.json               # semantic-release configuration
└── docker-compose.yaml           # local development
```


## CI/CD Pipeline

The project includes a modular GitHub Actions pipeline split into separate workflows for better maintainability:

```
.github/workflows/
├── ci-test.yaml         # Testing for all microservices
├── ci-build.yaml        # Build and push Docker images
├── ci-security.yaml     # Trivy security scanning
├── ci-helm.yaml         # Helm chart validation
├── release.yaml         # Semantic release automation
├── renovate.yaml        # Automated dependency updates
└── pipeline-summary.yaml # Aggregated status
```

### Features

- **Change detection**: Only builds services that have changed
- **Matrix builds**: Parallel testing and building for all services
- **Multi-language support**: Go, Java, C#, Ruby, and Node.js testing
- **GHCR integration**: Automatic image publishing to GitHub Container Registry
- **Security scanning**: Trivy vulnerability scanning for all images and IaC
- **Helm validation**: Chart linting and template verification
- **Semantic Release**: Automated versioning based on Conventional Commits
- **Renovate**: Automated dependency updates with automerge for patch/minor

### Triggering Builds

- **Push to main**: Full build, push to GHCR, and potential release
- **Pull requests**: Test and lint only
- **Manual dispatch**: Force build all services

### Image Tags

Images are published to GHCR with the following tags:
- `v{x.y.z}` - Semantic version (created automatically)
- `sha-<commit>` - Git commit SHA
- `latest` - Latest main branch build

```bash
# pull images
docker pull ghcr.io/sabinghosty19/test-workflow/api-gateway:latest
docker pull ghcr.io/sabinghosty19/test-workflow/order-service:latest
docker pull ghcr.io/sabinghosty19/test-workflow/inventory-service:latest
docker pull ghcr.io/sabinghosty19/test-workflow/notification-service:latest
docker pull ghcr.io/sabinghosty19/test-workflow/frontend:latest
```

### Semantic Versioning

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning:

| Commit Type | Version Bump | Example |
|-------------|--------------|---------|
| `fix:` | Patch (0.0.X) | `fix: resolve null pointer` |
| `feat:` | Minor (0.X.0) | `feat: add new endpoint` |
| `feat!:` | Major (X.0.0) | `feat!: change API format` |

See [.github/docs/SEMANTIC_RELEASE.md](.github/docs/SEMANTIC_RELEASE.md) for complete guide.

### Security Scanning

Trivy scans are run on every build and produce HTML reports available in GitHub Artifacts. Reports also appear in the GitHub Security tab.

See [.github/docs/TRIVY.md](.github/docs/TRIVY.md) for more information.

### Dependency Updates

Renovate automatically creates PRs for dependency updates:
- **Daily schedule**: Checks for updates every weekday
- **Automerge**: Patch and minor updates are merged automatically
- **Grouped PRs**: Updates are grouped by language

See [.github/docs/RENOVATE.md](.github/docs/RENOVATE.md) for configuration details.

## Helm Chart

The Helm chart provides a modular deployment of all microservices.

### Installation

```bash
# update dependencies
helm dependency update helm/test-workflow

# install with default values
helm install test-workflow helm/test-workflow -n test-workflow --create-namespace

# install with custom values
helm install test-workflow helm/test-workflow -n test-workflow \
  --set global.imageTag=v1.0 \
  --set ingress.hosts[0].host=myapp.example.com
```

### Configuration

Key configuration options in `values.yaml`:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.imageRegistry` | Container registry | `ghcr.io/sabinghosty19/test-workflow` |
| `global.imageTag` | Image tag for all services | `v1.0` |
| `apiGateway.replicaCount` | API Gateway replicas | `2` |
| `orderService.replicaCount` | Order Service replicas | `2` |
| `ingress.enabled` | Enable ingress | `true` |
| `postgresql.enabled` | Deploy PostgreSQL | `true` |
| `kafka.enabled` | Deploy Kafka | `true` |
| `networkPolicies.enabled` | Enable network policies | `true` |

## API Endpoints

### Orders

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/orders` | Create a new order |
| GET | `/api/v1/orders` | List all orders |
| GET | `/api/v1/orders/:id` | Get order details |
| PUT | `/api/v1/orders/:id/status` | Update order status |

### Inventory

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/inventory/products` | List products |
| POST | `/api/v1/inventory/products` | Add a product |
| GET | `/api/v1/inventory/products/:id` | Get product stock |
| PUT | `/api/v1/inventory/products/:id/stock` | Update stock |

### Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/notifications` | Get notification history |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |

## Zero Trust Security

The project implements zero trust security principles through Kubernetes Network Policies:

### Network Policies

- **Default Deny**: All ingress traffic is blocked by default
- **Frontend**: Allows ingress from anywhere (public-facing)
- **API Gateway**: Allows ingress from frontend and ingress controller
- **Order Service**: Allows gRPC only from API Gateway
- **Inventory Service**: Allows gRPC from API Gateway and Order Service
- **Notification Service**: Allows gRPC from API Gateway and Order Service
- **PostgreSQL**: Allows connections only from backend tier pods
- **Kafka**: Allows connections only from backend tier pods

### Security Features

- Pod Security Standards enforced (restricted mode in production)
- Read-only root filesystems where possible
- Non-root container execution
- Capability dropping
- Resource limits on all containers

## Testing

```bash
# api gateway (go)
cd services/api-gateway && go test ./...

# order service (java)
cd services/order-service && mvn test

# inventory service (c#)
cd services/inventory-service && dotnet test

# notification service (ruby)
cd services/notification-service && bundle exec rspec

# frontend (react)
cd frontend && npm test
```

## Environment Variables

### Common Variables

```bash
DB_HOST=postgres
DB_PORT=5432
DB_NAME=testworkflow
DB_USER=postgres
DB_PASSWORD=postgres
KAFKA_BOOTSTRAP_SERVERS=kafka:9092
ORDER_SERVICE_ADDR=order-service:50051
INVENTORY_SERVICE_ADDR=inventory-service:50052
NOTIFICATION_SERVICE_ADDR=notification-service:50053
```

### Service-Specific

| Service | Variable | Description |
|---------|----------|-------------|
| API Gateway | `GIN_MODE` | Gin framework mode (release/debug) |
| Order Service | `JAVA_OPTS` | JVM options |
| Order Service | `SPRING_PROFILES_ACTIVE` | Active Spring profiles |
| Inventory Service | `ASPNETCORE_URLS` | ASP.NET Core URLs |

## Local Ports

| Service | Port |
|---------|------|
| Frontend | 3000 |
| API Gateway | 8080 |
| Order Service | 8081, 50051 |
| Inventory Service | 8082, 50052 |
| Notification Service | 8083, 50053 |
| PostgreSQL | 5432 |
| Kafka | 9092 |
| Zookeeper | 2181 |

## Flux CD Environments

The project includes Flux CD configurations for three environments:

| Environment | Path | Description |
|-------------|------|-------------|
| Development | `flux/overlays/development` | Minimal resources, latest tags, no network policies |
| Staging | `flux/overlays/staging` | Moderate resources, version tags, network policies enabled |
| Production | `flux/overlays/production` | High availability, autoscaling, TLS, strict security |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests locally
5. Submit a pull request

## License

MIT License - see LICENSE file for details
