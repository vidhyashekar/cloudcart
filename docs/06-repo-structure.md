# Repository Structure

Each of Auth, Product, Order, and Notification should be an independent microservice.

The fact that they are under one repository does not make them one application. The structure I suggested is a monorepo.

## One GitHub repository — Monorepo
```
Detailed structure:

cloudcart/
│
├── services/
│   ├── auth-service/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── go.mod
│   │   └── Dockerfile
│   │
│   ├── product-service/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── go.mod
│   │   └── Dockerfile
│   │
│   ├── order-service/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── go.mod
│   │   └── Dockerfile
│   │
│   └── notification-service/
│       ├── cmd/
│       ├── internal/
│       ├── go.mod
│       └── Dockerfile
│
├── deployments/
│   ├── docker/
│   │   └── docker-compose.yml
│   │
│   ├── kubernetes/
│   │   ├── auth/
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   ├── ingress.yaml
│   │   ├── product/
│   │   ├── order/
│   │   └── notification/
│   │
│   └── helm/
│
├── infrastructure/
│   └── terraform/
│   └── ansible/
│
├── docs/
│
└── .github/
    └── workflows/
        ├── auth-service.yml
        ├── product-service.yml
        ├── order-service.yml
        ├── notification-service.yml
        └── infrastructure.yml
```

```
Compressed structure
cloudcart/
│
├── services/
│   ├── auth-service/
│   ├── product-service/
│   ├── order-service/
│   └── notification-service/
│
├── docs/
├── deployments/
├── infrastructure/
└── .github/
```

Each service:

- Has its own go.mod
- Has its own Dockerfile
- Builds its own Docker image
- Is deployed independently to Kubernetes
- Can be scaled independently
- Can be released independently

Each pipeline can:

- Build only that service
- Run only that service's tests
- Build only that Docker image
- Push only that Docker image
- Deploy only that Kubernetes Deployment

Every service has its own Dockerfile and each image is built independently.
```
docker build -t auth-service .
docker build -t product-service .
docker build -t order-service .
```

Exactly like separate repositories.

For example:
```
auth-service       → cloudcart/auth-service:v1.0.0
product-service    → cloudcart/product-service:v1.0.0
order-service      → cloudcart/order-service:v1.0.0
notification-service → cloudcart/notification-service:v1.0.0
```
Every service has its own manifests.

This gives you the best learning experience because we can see the entire system in one place, while still implementing every service as an independent microservice.

```
                    GitHub Repository
                         cloudcart
                             │
       ┌─────────────────────┼─────────────────────┐
       │                     │                     │
 Auth Service        Product Service        Order Service
       │                     │                     │
  Docker Image        Docker Image          Docker Image
       │                     │                     │
  Kubernetes Pod      Kubernetes Pod        Kubernetes Pod
```

Each service can have its own CI/CD pipeline logic as well.
We should build four individual microservices; we are only keeping them inside one repository to make the project easier to manage and learn from.

## Summary
✅ One Git repository (monorepo) for easier management.
✅ Four independent Go modules (one per service).
✅ One Dockerfile per service.
✅ Separate GitHub Actions workflow per service.
✅ Separate Kubernetes manifests (or Helm charts) per service.
✅ A single docker-compose.yml for local development that starts the entire CloudCart ecosystem.