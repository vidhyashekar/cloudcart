# cloudcart
CloudCart is a cloud-native e-commerce backend application designed to demonstrate production-grade software architecture, cloud-native development, DevOps practices, and microservices.

## Architecture

> Refer: `docs/` directory for the complete architecture.

## Docker Networking

                    Docker Network
                         │
        ┌────────────────┼────────────────┐
        │                │                │
     Auth:8080       Product:8081      Order:8082
        │                ▲                │
        │                │                │
        └───────┐        │        ┌───────┘
                ▼        │        ▼
             PostgreSQL  │      Kafka:9092
                         │         │
                         │         ▼
                         │   Notification
                         │
                  HTTP service-to-service