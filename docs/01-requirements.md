# CloudCart - Requirements Analysis

## Document Information

| Field | Value |
|-------|-------|
| Project | CloudCart |
| Version | 1.0 |
| Author | Vidhya Shekar |
| Status | Draft |
| Last Updated | July 2026 |

---

# 1. Project Overview

## Purpose

CloudCart is a cloud-native e-commerce backend application designed to demonstrate production-grade software architecture, cloud-native development, DevOps practices, and microservices.

The primary objective of this project is to build a scalable, maintainable, and highly available backend system using Go, Kubernetes, Docker, AWS, and modern DevOps tooling.

This project serves as:

- A portfolio project
- A hands-on DevOps learning project
- A production-grade backend reference implementation
- Technical Architect interview preparation

---

# 2. Business Problem

Traditional monolithic applications become difficult to maintain as systems grow.

Some common challenges include:

- Large codebase
- Difficult deployments
- Limited scalability
- Tight coupling
- Single point of failure
- Longer release cycles

CloudCart addresses these challenges using a Microservices Architecture.

---

# 3. Goals

The application should:

- Allow customers to browse products.
- Allow customers to register and login.
- Allow customers to place orders.
- Send notifications after successful order creation.
- Support independent service deployments.
- Support horizontal scaling.
- Demonstrate event-driven communication.
- Be cloud-ready.
- Follow production-grade architecture.

---

# 4. Scope

## In Scope

✔ User Authentication

✔ Product Management

✔ Order Management

✔ Notification Service

✔ REST APIs

✔ JWT Authentication

✔ Docker

✔ Kubernetes

✔ CI/CD

✔ Monitoring

✔ Logging

✔ AWS Deployment

---

## Out of Scope

The following features will not be implemented in Version 1.

- Payment Gateway
- Shopping Cart
- Wishlist
- Reviews & Ratings
- Product Recommendation Engine
- Search Engine (ElasticSearch)
- Mobile Application
- Admin Dashboard UI

These can be added in future releases.

---

# 5. Functional Requirements

## Authentication

The system shall allow users to:

- Register
- Login
- Receive JWT tokens
- Access protected APIs

---

## Product Management

The system shall allow users to:

- View products
- View product details
- Search products

Future versions may support:

- Create Product
- Update Product
- Delete Product

---

## Order Management

The system shall allow users to:

- Create Order
- View Orders
- View Order Details

---

## Notification

The system shall:

- Send Order Confirmation
- Send Welcome Email
- Support future notification channels

---

# 6. Non-Functional Requirements

## Availability

Target Availability:

99.9%

---

## Scalability

The application shall support:

- Horizontal Scaling
- Independent Service Scaling

---

## Reliability

Failures in one microservice should not impact other services.

---

## Performance

Average API response time:

< 300 ms

---

## Security

The application shall support:

- JWT Authentication
- HTTPS
- Password Hashing
- Secure API Access

---

## Maintainability

The system should support:

- Independent deployment
- Loose coupling
- Clean Architecture
- Modular code

---

# 7. Users

## Customer

Responsibilities:

- Register
- Login
- Browse Products
- Place Orders

---

## Administrator

Responsibilities:

- Manage Products
- View Orders
- Monitor Application

---

# 8. Architecture Decision

The application will follow a Microservices Architecture.

Reasons include:

- Independent deployment
- Independent scaling
- Better fault isolation
- Clear ownership
- Faster development
- Better maintainability

The architecture choice is based on project requirements and learning objectives.

---

# 9. Technology Stack

| Layer | Technology |
|---------|------------|
| Language | Go |
| Framework | Gin |
| Database | PostgreSQL |
| Authentication | JWT |
| Containerization | Docker |
| Orchestration | Kubernetes |
| Messaging | Kafka |
| CI/CD | GitHub Actions |
| GitOps | Argo CD |
| Monitoring | Prometheus |
| Dashboard | Grafana |
| Cloud | AWS |

---

# 10. Assumptions

- Services are stateless.
- Each service owns its own database.
- Communication is REST unless asynchronous processing is required.
- Kafka is used for event-driven communication.
- Kubernetes manages deployments.
- AWS is the target cloud provider.

---

# 11. Risks

| Risk | Mitigation |
|------|------------|
| Service communication complexity | API Gateway and well-defined APIs |
| Event processing failures | Kafka retries and consumer groups |
| Scaling issues | Kubernetes Horizontal Pod Autoscaler |
| Deployment failures | CI/CD pipeline with automated rollback |

---

# 12. Future Enhancements

Future releases may include:

- Payment Service
- Shopping Cart
- Redis Cache
- Elasticsearch
- Service Mesh (Istio)
- Distributed Tracing
- Circuit Breaker Pattern
- CQRS
- Event Sourcing
- Multi-region Deployment

---

# 13. Related Documents

- 02-high-level-design.md
- 03-service-boundaries.md
- 04-database-design.md
- 05-api-design.md