# Test Workflow - Zero Trust Microservices Demo

This project demonstrates zero trust security patterns in Kubernetes using a polyglot microservices architecture. The application simulates an e-commerce order management system where each service is written in a different programming language and communicates via REST and gRPC protocols.

The primary goal is to showcase how network policies, service mesh principles, and least-privilege access controls can be implemented in a Kubernetes environment to achieve defense in depth.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [System Components](#system-components)
4. [Communication Patterns](#communication-patterns)
5. [Zero Trust Implementation](#zero-trust-implementation)
6. [Getting Started](#getting-started)
7. [Kubernetes Deployment](#kubernetes-deployment)
8. [API Reference](#api-reference)
9. [CI/CD Pipeline](#cicd-pipeline)
10. [Configuration Reference](#configuration-reference)
11. [Troubleshooting](#troubleshooting)
12. [References](#references)

---

## Project Overview

Modern cloud-native applications require security models that assume no implicit trust between services. This demo implements such a model through:

- Explicit service-to-service authentication via gRPC with mTLS capability
- Network segmentation using Kubernetes NetworkPolicies
- Least-privilege access where each service can only communicate with its designated peers
- Defense in depth with multiple security layers

The application itself is a simplified order management system with inventory tracking and notification capabilities.

---

## Architecture

### High-Level System Diagram

```
+------------------------------------------------------------------+
|                                                                  |
|                     INGRESS CONTROLLER                           |
|                          (nginx)                                 |
|                                                                  |
+----------------------------------+-------------------------------+
                                   |
                                   | HTTP/443
                                   v
+------------------------------------------------------------------+
|                                                                  |
|                    FRONTEND (React + Nginx)                      |
|                         Port: 80                                 |
|                                                                  |
+----------------------------------+-------------------------------+
                                   |
                                   | REST/HTTP
                                   v
+------------------------------------------------------------------+
|                                                                  |
|                    API GATEWAY (Go + Gin)                        |
|                        Port: 8080                                |
|    - Request routing and aggregation                             |
|    - Protocol translation (REST to gRPC)                         |
|    - Rate limiting and request validation                        |
|                                                                  |
+----------+-------------------+-------------------+---------------+
           |                   |                   |
           | gRPC/50051        | gRPC/50052        | gRPC/50053
           v                   v                   v
+----------------+    +------------------+    +--------------------+
|                |    |                  |    |                    |
| ORDER SERVICE  |    | INVENTORY        |    | NOTIFICATION       |
| (Java/Spring)  |    | SERVICE          |    | SERVICE            |
|                |    | (C#/.NET)        |    | (Go)               |
| - Order CRUD   |    |                  |    |                    |
| - Status mgmt  |    | - Stock mgmt     |    | - Order alerts     |
| - Kafka prod.  |    | - Reservations   |    | - Stock alerts     |
|                |    | - Kafka cons.    |    | - Email dispatch   |
+-------+--------+    +--------+---------+    +---------+----------+
        |                      |                        |
        |                      |                        |
        +----------------------+------------------------+
                               |
                               v
              +--------------------------------+
              |                                |
              |     POSTGRESQL DATABASE        |
              |         Port: 5432             |
              |    (Schema per service)        |
              |                                |
              +--------------------------------+

              +--------------------------------+
              |                                |
              |      APACHE KAFKA              |
              |        Port: 9092              |
              |   (Event-driven messaging)     |
              |                                |
              +--------------------------------+
```

### Service Communication Flow

```
                            Request Flow
                            ============

  Client                                                    Database
    |                                                          |
    |  1. HTTP POST /api/v1/orders                            |
    +-------------------------->                               |
    |                          |                               |
    |                    API Gateway                           |
    |                          |                               |
    |    2. gRPC CreateOrder() |                               |
    |                          +-------------->                |
    |                          |              |                |
    |                          |        Order Service          |
    |                          |              |                |
    |                          |   3. Check inventory          |
    |                          |              +--------------->|
    |                          |              |                |
    |                          |        Inventory Service      |
    |                          |              |                |
    |                          |   4. Reserve stock            |
    |                          |              +--------------->|
    |                          |              |           PostgreSQL
    |                          |              |                |
    |                          |   5. Kafka event              |
    |                          |              +------+         |
    |                          |              |      |         |
    |                          |              |    Kafka       |
    |                          |              |      |         |
    |                          |   6. Consume event            |
    |                          |              |<-----+         |
    |                          |              |                |
    |                          |        Notification Service   |
    |                          |              |                |
    |                          |   7. Send notification        |
    |                          |              +--------------->|
    |                          |                          PostgreSQL
    |  8. HTTP 201 Created     |                               |
    <--------------------------+                               |
    |                                                          |
```

### Network Policy Diagram

```
                    Zero Trust Network Segmentation
                    ================================

    +------------------------------------------------------------------+
    |                      KUBERNETES CLUSTER                          |
    |                                                                  |
    |   NAMESPACE: test-workflow                                       |
    |                                                                  |
    |   +------------------+        DEFAULT DENY ALL INGRESS           |
    |   |    FRONTEND      |<---+                                      |
    |   |   (public tier)  |    |   Allowed: ingress-nginx             |
    |   +--------+---------+    |                                      |
    |            |              |                                      |
    |            | allowed      |                                      |
    |            v              |                                      |
    |   +------------------+    |                                      |
    |   |   API GATEWAY    |<---+   Allowed: frontend, ingress-nginx   |
    |   |  (gateway tier)  |                                           |
    |   +--------+---------+                                           |
    |            |                                                     |
    |            | allowed (gRPC only)                                 |
    |            v                                                     |
    |   +------------------+------------------+------------------+     |
    |   |                  |                  |                  |     |
    |   | ORDER SERVICE    | INVENTORY SVC    | NOTIFICATION SVC |     |
    |   | (backend tier)   | (backend tier)   | (backend tier)   |     |
    |   |                  |                  |                  |     |
    |   | Allowed from:    | Allowed from:    | Allowed from:    |     |
    |   | - api-gateway    | - api-gateway    | - api-gateway    |     |
    |   |                  | - order-service  | - order-service  |     |
    |   +--------+---------+--------+---------+--------+---------+     |
    |            |                  |                  |               |
    |            +------------------+------------------+               |
    |                               |                                  |
    |                               | allowed (port 5432)              |
    |                               v                                  |
    |                      +------------------+                        |
    |                      |   POSTGRESQL     |                        |
    |                      |   (data tier)    |                        |
    |                      |                  |                        |
    |                      | Allowed from:    |                        |
    |                      | - backend tier   |                        |
    |                      +------------------+                        |
    |                                                                  |
    +------------------------------------------------------------------+
```

---

## System Components

### Service Overview

| Service | Language | Framework | HTTP Port | gRPC Port | Purpose |
|---------|----------|-----------|-----------|-----------|---------|
| Frontend | TypeScript | React 18 | 80 | - | Single-page application served via Nginx |
| API Gateway | Go 1.21 | Gin | 8080 | - | REST entry point, protocol translation |
| Order Service | Java 21 | Spring Boot 3 | 8081 | 50051 | Order lifecycle management |
| Inventory Service | C# | .NET 8 | 8082 | 50052 | Stock and reservation handling |
| Notification Service | Go 1.21 | grpc-go | 8083 | 50053 | Alert and notification dispatch |

### Infrastructure Components

| Component | Version | Port | Purpose |
|-----------|---------|------|---------|
| PostgreSQL | 15 | 5432 | Persistent storage with schema-per-service |
| Apache Kafka | 3.6 | 9092 | Asynchronous event messaging |
| Zookeeper | 3.8 | 2181 | Kafka coordination |
| Nginx Ingress | 1.9 | 80/443 | External traffic routing |

---

## Communication Patterns

### REST API (External)

External clients interact with the system through the API Gateway using standard REST conventions:

```
POST   /api/v1/orders           Create new order
GET    /api/v1/orders           List orders (paginated)
GET    /api/v1/orders/{id}      Get order details
PUT    /api/v1/orders/{id}      Update order status
```

### gRPC (Internal)

Internal service communication uses Protocol Buffers for efficient serialization:

```protobuf
service OrderService {
    rpc CreateOrder(CreateOrderRequest) returns (OrderResponse);
    rpc GetOrder(GetOrderRequest) returns (OrderResponse);
    rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
    rpc UpdateOrderStatus(UpdateStatusRequest) returns (StatusResponse);
}

service InventoryService {
    rpc GetStock(GetStockRequest) returns (StockResponse);
    rpc ReserveStock(ReserveStockRequest) returns (StatusResponse);
    rpc ReleaseStock(ReleaseStockRequest) returns (StatusResponse);
}

service NotificationService {
    rpc SendOrderNotification(OrderNotificationRequest) returns (StatusResponse);
    rpc SendStockAlert(StockAlertRequest) returns (StatusResponse);
    rpc GetNotifications(GetNotificationsRequest) returns (NotificationsResponse);
}
```

### Event-Driven (Kafka)

Asynchronous communication between services:

| Topic | Producer | Consumer | Payload |
|-------|----------|----------|---------|
| order-events | Order Service | Inventory, Notification | Order state changes |
| inventory-events | Inventory Service | Notification | Stock level alerts |

---

## Zero Trust Implementation

### Principles Applied

1. **Never Trust, Always Verify**: Every service request is authenticated and authorized
2. **Least Privilege Access**: Services can only reach their designated dependencies
3. **Assume Breach**: Network segmentation limits blast radius
4. **Explicit Allow**: Default-deny policies with explicit allowlists

### Network Policies

The Helm chart includes predefined NetworkPolicy resources:

```yaml
# Default deny all ingress
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
spec:
  podSelector: {}
  policyTypes:
    - Ingress

# Allow API Gateway to reach Order Service
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-api-gateway-to-order
spec:
  podSelector:
    matchLabels:
      app: order-service
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: api-gateway
      ports:
        - protocol: TCP
          port: 50051
```

### Security Controls

| Control | Implementation |
|---------|----------------|
| Pod Security Standards | Restricted mode in production |
| Container Security | Non-root user, read-only filesystem |
| Capability Dropping | All capabilities dropped except required |
| Resource Limits | CPU and memory limits on all pods |
| Secret Management | Kubernetes Secrets with encryption at rest |

---

## Getting Started

### Prerequisites

- Docker 24.x or later
- Docker Compose 2.x
- Kubernetes 1.28+ (minikube, kind, or cloud provider)
- Helm 3.13+
- kubectl 1.28+
- Go 1.21+ (for local development)

### Local Development with Docker Compose

```bash
# Clone the repository
git clone https://github.com/sabinghosty19/test-workflow.git
cd test-workflow

# Build all service images
./scripts/build-all.sh

# Start the full stack
docker-compose up -d

# Verify services are running
docker-compose ps

# View logs
docker-compose logs -f api-gateway

# Access the application
# Frontend: http://localhost:3000
# API Gateway: http://localhost:8080

# Stop all services
docker-compose down
```

### Running Tests

```bash
# API Gateway (Go)
cd services/api-gateway && go test -v ./...

# Order Service (Java)
cd services/order-service && mvn test

# Inventory Service (C#)
cd services/inventory-service && dotnet test

# Notification Service (Go)
cd services/notification-service && go test -v ./...

# Frontend (React)
cd frontend && npm test
```

---

## Kubernetes Deployment

### Using Helm

```bash
# Add Helm repository dependencies
cd helm/test-workflow
helm dependency update

# Create namespace
kubectl create namespace test-workflow

# Install with default values
helm install test-workflow ./helm/test-workflow \
  --namespace test-workflow

# Install with custom values
helm install test-workflow ./helm/test-workflow \
  --namespace test-workflow \
  --set global.imageTag=v1.2.0 \
  --set ingress.hosts[0].host=app.example.com \
  --set networkPolicies.enabled=true

# Verify deployment
kubectl get pods -n test-workflow
kubectl get svc -n test-workflow

# Check service health
kubectl exec -it deploy/api-gateway -n test-workflow -- wget -qO- http://localhost:8080/health
```

### Using Kustomize

```bash
# Deploy base configuration
kubectl apply -k k8s/base

# Deploy with overlays
kubectl apply -k k8s/overlays/production

# Add local DNS entry for testing
echo "127.0.0.1 test-workflow.local" | sudo tee -a /etc/hosts
```

### GitOps with Flux CD

```bash
# Bootstrap Flux
flux bootstrap github \
  --owner=sabinghosty19 \
  --repository=test-workflow \
  --path=flux/overlays/development \
  --personal

# Monitor reconciliation
flux get helmreleases -n flux-system
flux get kustomizations -A

# Force reconciliation
flux reconcile helmrelease test-workflow -n flux-system
```

### Environment Configurations

| Environment | Path | Characteristics |
|-------------|------|-----------------|
| Development | flux/overlays/development | Minimal resources, debug logging, no network policies |
| Staging | flux/overlays/staging | Moderate resources, version tags, network policies enabled |
| Production | flux/overlays/production | HA replicas, autoscaling, TLS, strict security |

---

## API Reference

### Orders API

#### Create Order

```http
POST /api/v1/orders
Content-Type: application/json

{
  "customer_name": "John Doe",
  "items": [
    {"product_id": "prod-123", "quantity": 2}
  ]
}
```

Response:
```json
{
  "id": "ord-456",
  "customer_name": "John Doe",
  "status": "PENDING",
  "total_amount": 59.99,
  "created_at": "2025-01-15T10:30:00Z"
}
```

#### List Orders

```http
GET /api/v1/orders?limit=10&offset=0
```

#### Get Order

```http
GET /api/v1/orders/{id}
```

#### Update Order Status

```http
PUT /api/v1/orders/{id}/status
Content-Type: application/json

{
  "status": "SHIPPED"
}
```

### Inventory API

#### List Products

```http
GET /api/v1/inventory/products
```

#### Get Product Stock

```http
GET /api/v1/inventory/products/{id}
```

#### Update Stock

```http
PUT /api/v1/inventory/products/{id}/stock
Content-Type: application/json

{
  "quantity": 100,
  "operation": "SET"
}
```

### Notifications API

#### Get Notifications

```http
GET /api/v1/notifications?limit=10&offset=0
```

### Health Endpoints

All services expose health endpoints:

| Service | Endpoint |
|---------|----------|
| API Gateway | GET /health |
| Order Service | GET /actuator/health |
| Inventory Service | GET /health |
| Notification Service | GET /health |

---

## CI/CD Pipeline

The project uses GitHub Actions with modular workflows. See [.github/README.md](.github/README.md) for complete pipeline documentation.

### Workflow Overview

```
                        CI/CD Pipeline Flow
                        ===================

  +-------------+     +-------------+     +----------------+
  |             |     |             |     |                |
  |  Developer  +---->+  Pull       +---->+  CI Tests      |
  |  Commit     |     |  Request    |     |  (all langs)   |
  |             |     |             |     |                |
  +-------------+     +------+------+     +-------+--------+
                             |                    |
                             | merge              | pass
                             v                    v
                      +------+------+     +-------+--------+
                      |             |     |                |
                      |  main       +---->+  Build & Push  |
                      |  branch     |     |  (GHCR)        |
                      |             |     |                |
                      +------+------+     +-------+--------+
                             |                    |
                             |                    v
                             |            +-------+--------+
                             |            |                |
                             +----------->+  Security Scan |
                                          |  (Trivy)       |
                                          |                |
                                          +-------+--------+
                                                  |
                                                  v
                                          +-------+--------+
                                          |                |
                                          |  Release       |
                                          |  (semantic)    |
                                          |                |
                                          +----------------+
```

### Available Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| ci-test.yaml | Push, PR | Run tests for all services |
| ci-build.yaml | Push to main | Build and push Docker images |
| ci-security.yaml | After build | Trivy vulnerability scanning |
| ci-helm.yaml | Push, PR | Helm chart validation |
| release.yaml | Push to main | Semantic versioning and release |
| renovate.yaml | Scheduled | Dependency updates |

### Docker Images

Images are published to GitHub Container Registry:

```
ghcr.io/sabinghosty19/test-workflow/api-gateway:latest
ghcr.io/sabinghosty19/test-workflow/order-service:latest
ghcr.io/sabinghosty19/test-workflow/inventory-service:latest
ghcr.io/sabinghosty19/test-workflow/notification-service:latest
ghcr.io/sabinghosty19/test-workflow/frontend:latest
```

Tag formats:
- `latest` - Latest main branch build
- `v1.2.3` - Semantic version release
- `sha-abc1234` - Specific commit

---

## Configuration Reference

### Environment Variables

#### Common Variables

```bash
DB_HOST=postgres
DB_PORT=5432
DB_NAME=testworkflow
DB_USER=postgres
DB_PASSWORD=postgres
KAFKA_BOOTSTRAP_SERVERS=kafka:9092
```

#### Service-Specific Variables

| Service | Variable | Default | Description |
|---------|----------|---------|-------------|
| API Gateway | GIN_MODE | release | Gin framework mode |
| API Gateway | ORDER_SERVICE_ADDR | order-service:50051 | gRPC endpoint |
| API Gateway | INVENTORY_SERVICE_ADDR | inventory-service:50052 | gRPC endpoint |
| API Gateway | NOTIFICATION_SERVICE_ADDR | notification-service:50053 | gRPC endpoint |
| Order Service | JAVA_OPTS | -Xmx512m | JVM options |
| Order Service | SPRING_PROFILES_ACTIVE | default | Spring profile |
| Inventory Service | ASPNETCORE_URLS | http://+:8082 | ASP.NET URLs |
| Notification Service | GRPC_PORT | 50053 | gRPC listen port |
| Notification Service | HTTP_PORT | 8083 | Health check port |

### Helm Values

Key configuration options in `values.yaml`:

| Parameter | Default | Description |
|-----------|---------|-------------|
| global.imageRegistry | ghcr.io/sabinghosty19/test-workflow | Container registry |
| global.imageTag | latest | Default image tag |
| apiGateway.replicaCount | 2 | API Gateway replicas |
| orderService.replicaCount | 2 | Order Service replicas |
| inventoryService.replicaCount | 2 | Inventory Service replicas |
| notificationService.replicaCount | 1 | Notification Service replicas |
| ingress.enabled | true | Enable ingress |
| ingress.className | nginx | Ingress class |
| postgresql.enabled | true | Deploy PostgreSQL |
| kafka.enabled | true | Deploy Kafka |
| networkPolicies.enabled | true | Enable network policies |

### Port Mapping

| Service | Container Port | Service Port | NodePort (dev) |
|---------|----------------|--------------|----------------|
| Frontend | 80 | 80 | 30080 |
| API Gateway | 8080 | 8080 | 30880 |
| Order Service | 50051 | 50051 | - |
| Inventory Service | 50052 | 50052 | - |
| Notification Service | 50053 | 50053 | - |
| PostgreSQL | 5432 | 5432 | - |
| Kafka | 9092 | 9092 | - |

---

## Troubleshooting

### Common Issues

#### Pods not starting

```bash
# Check pod status
kubectl get pods -n test-workflow

# View pod events
kubectl describe pod <pod-name> -n test-workflow

# Check logs
kubectl logs <pod-name> -n test-workflow
```

#### Service connectivity issues

```bash
# Test gRPC connectivity from API Gateway
kubectl exec -it deploy/api-gateway -n test-workflow -- \
  grpcurl -plaintext order-service:50051 list

# Check network policies
kubectl get networkpolicies -n test-workflow

# Verify service endpoints
kubectl get endpoints -n test-workflow
```

#### Database connection failures

```bash
# Check PostgreSQL status
kubectl get pods -l app=postgresql -n test-workflow

# Test database connectivity
kubectl exec -it deploy/api-gateway -n test-workflow -- \
  pg_isready -h postgres -p 5432 -U postgres
```

### Debug Commands

```bash
# Port forward for local debugging
kubectl port-forward svc/api-gateway 8080:8080 -n test-workflow

# View all resources
kubectl get all -n test-workflow

# Check Helm release status
helm status test-workflow -n test-workflow

# View Helm values
helm get values test-workflow -n test-workflow
```

---

## References

### Kubernetes and Cloud Native

- Kubernetes Documentation: https://kubernetes.io/docs/
- Kubernetes Network Policies: https://kubernetes.io/docs/concepts/services-networking/network-policies/
- Helm Documentation: https://helm.sh/docs/
- Flux CD Documentation: https://fluxcd.io/docs/

### Zero Trust Security

- NIST Zero Trust Architecture (SP 800-207): https://csrc.nist.gov/publications/detail/sp/800-207/final
- CNCF Zero Trust Whitepaper: https://www.cncf.io/blog/2021/11/18/zero-trust-networking/
- Kubernetes Security Best Practices: https://kubernetes.io/docs/concepts/security/

### gRPC and Protocol Buffers

- gRPC Documentation: https://grpc.io/docs/
- Protocol Buffers Guide: https://protobuf.dev/programming-guides/proto3/
- gRPC-Go: https://github.com/grpc/grpc-go

### Container Security

- Trivy Scanner: https://aquasecurity.github.io/trivy/
- Docker Security Best Practices: https://docs.docker.com/develop/security-best-practices/
- Pod Security Standards: https://kubernetes.io/docs/concepts/security/pod-security-standards/

### CI/CD

- GitHub Actions: https://docs.github.com/en/actions
- Semantic Release: https://semantic-release.gitbook.io/
- Conventional Commits: https://www.conventionalcommits.org/

### Programming Languages and Frameworks

- Go: https://go.dev/doc/
- Gin Web Framework: https://gin-gonic.com/docs/
- Spring Boot: https://docs.spring.io/spring-boot/
- ASP.NET Core: https://docs.microsoft.com/en-us/aspnet/core/
- React: https://react.dev/

---

## License

MIT License - see LICENSE file for details.
