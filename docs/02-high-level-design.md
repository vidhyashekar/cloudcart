# High Level Design

## Overview

CloudCart is a cloud-native e-commerce application built using a microservices architecture. The system is designed to be scalable, highly available, loosely coupled, and production-ready. Each business capability is implemented as an independent Go microservice and deployed on Amazon EKS.

The architecture uses synchronous REST communication for request-response operations and asynchronous event-driven communication using Kafka for background processing.

---

## High Level Architecture

> Refer: `docs/images/hld.png`

The system consists of the following major components:

- Application Load Balancer (ALB)
- API Gateway
- Authentication Service
- Product Service
- Order Service
- Notification Service
- Kafka Event Bus
- Amazon RDS PostgreSQL
- Amazon ElastiCache (Redis)
- Prometheus & Grafana

---

## Communication Pattern

### Synchronous Communication (REST)

- Client → API Gateway
- API Gateway → Authentication Service
- API Gateway → Product Service
- API Gateway → Order Service

Used whenever an immediate response is required.

### Asynchronous Communication (Kafka)

Example:

```
Order Created
        ↓
Kafka Topic (order.created)
        ↓
Notification Service
        ↓
Send Email
```

This ensures background operations do not block user requests.

---
## Authentication Strategy

Authentication is implemented using JWT.

Workflow:

- User logs in
- Authentication Service validates credentials
- JWT Access Token is generated
- API Gateway validates JWT
- Request is forwarded to backend services

## Data Ownership

Each microservice owns its own database.

| Service | Database |
|----------|----------|
| Authentication Service | PostgreSQL |
| Product Service | PostgreSQL |
| Order Service | PostgreSQL |
| Notification Service | PostgreSQL |

No service directly accesses another service's database.

---

## Architecture Principles

- Microservices architecture
- Database per service
- Loose coupling
- High availability
- Horizontal scalability
- Event-driven communication
- Stateless services
- Independent deployments

---

## Architecture Decisions

### Microservices

CloudCart uses a microservices architecture to enable independent deployment, better scalability, loose coupling, and clear separation of business capabilities.

### AWS

AWS is selected as the cloud provider because it offers managed services such as Amazon EKS, Amazon RDS, Amazon ElastiCache, and Amazon ECR, reducing operational overhead while providing production-grade reliability.

### Kubernetes

Amazon EKS is used to orchestrate containerized workloads, enabling self-healing, rolling deployments, autoscaling, and simplified application management.

### Database

Amazon RDS PostgreSQL is used instead of self-managed PostgreSQL on Kubernetes to reduce operational complexity and align with common production practices.

### Cache

Amazon ElastiCache (Redis) is used as the distributed cache to improve application performance and reduce database load.

### Messaging

Kafka is deployed as a self-managed service on Kubernetes (KRaft mode) to provide hands-on experience with deploying and operating distributed messaging systems.

### Networking

CloudCart uses two Availability Zones for high availability.

Each Availability Zone contains:

- One Public Subnet
- One Private Subnet

This simplified networking model provides a clean architecture suitable for learning while following AWS best practices. The design can later be extended with dedicated application, database, and cache subnets as the system grows.

---

## Future Improvements

- Horizontal Pod Autoscaler (HPA)
- Service Mesh (Istio)
- AWS WAF
- Secrets Manager
- Multi-region deployment
- Blue-Green deployment
- Canary deployment
- Distributed tracing using Jaeger



