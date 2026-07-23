# Database Design

## Database Strategy

CloudCart uses a single Amazon RDS PostgreSQL instance for simplicity and cost optimization.

The database is logically divided into separate schemas:

- auth
- product
- orders
- notification

This approach provides:

- Loose coupling
- Independent deployments
- Independent scaling
- Fault isolation
- Better maintainability

Each microservice exclusively owns its schema and is responsible for managing its own tables.

Services communicate through REST APIs and Kafka events. Direct access to another service's schema is not allowed, even though they reside in the same PostgreSQL database.

This design balances microservice principles with operational simplicity and is well suited for learning, portfolio projects, and small to medium-scale production systems.

## Database Selection

| Service | Database |
|----------|----------|
| Authentication Service | PostgreSQL |
| Product Service | PostgreSQL |
| Order Service | PostgreSQL |
| Notification Service | PostgreSQL |

Deployment:
- Amazon RDS PostgreSQL
- Multi-AZ enabled
- Automated backups
- Read replicas can be added later

```
Question: Why not MongoDB?
Because our domain (users, orders, order items, products) is highly relational. PostgreSQL provides ACID transactions, foreign keys, joins, and strong consistency, which are ideal for an e-commerce application.
```

## Database Ownership

```
Authentication Service
        │
        ▼
 Authentication Schema

Product Service
        │
        ▼
 Product Schema

Order Service
        │
        ▼
 Order Schema

Notification Service
        │
        ▼
 Notification Schema
```


Why this decision we could have gone for separate databases too?

In production, both approaches can be used:

Option 1 — Separate Databases
```
Amazon RDS Instance

auth_db
    ├── users
    ├── roles
    └── refresh_tokens

product_db
    ├── products
    ├── categories
    └── inventory

order_db
    ├── orders
    ├── order_items
    └── payments

notification_db
    └── notification_logs
```

This is common in large enterprises where each microservice is fully isolated.

Advantages
- Complete isolation
- Independent backups
- Independent scaling
- Easier to migrate one service

Disadvantages
- Higher AWS cost
- More operational overhead
- More databases to maintain

Option 2 — One Database, Multiple Schemas (Same PostgreSQL Instance)
```
Amazon RDS PostgreSQL

cloudcart
│
├── auth
│      ├── users
│      ├── roles
│      └── refresh_tokens
│
├── product
│      ├── products
│      ├── categories
│      └── inventory
│
├── orders
│      ├── orders
│      ├── order_items
│      └── payments
│
└── notification
       └── notification_logs
```
Notice:

- **One PostgreSQL database**
- Four schemas
- Each schema owns its own tables

This is still a valid implementation of the Database per Service pattern because ownership is enforced at the schema level.

### For this project

After having multiple thoughts recommending **one Amazon RDS PostgreSQL** instance with separate schemas:

```
cloudcart
│
├── auth
├── product
├── order
└── notification
```

Why?
- Easier to manage
- Lower AWS cost
- Simple backups
- Still respects service ownership
- Perfect for a portfolio project

## Tables

### Authentication Schema

```
- users
- roles
- refresh_tokens
```

### Product Schema
```
- products
- categories
- inventory
```

### Order Schema
```
- orders
- order_items
- payments
```

### Notification Schema
```
- Notification Schema
```

## Relationships

```
User

1 ------ N

Orders

Orders

1 ------ N

Order Items

Category

1 ------ N

Products

Products

1 ------ 1

Inventory
```

## Indexing Strategy

```
Primary Keys

UUID

Indexes

users.email

products.name

orders.user_id

orders.created_at

order_items.order_id

notification_logs.created_at
```
## Service Ownership
```
Authentication Service
        │
        ▼
auth schema
        │
users
roles
refresh_tokens

Product Service
        │
        ▼
product schema
        │
products
categories
inventory

Order Service
        │
        ▼
orders schema
        │
orders
order_items
payments

Notification Service
        │
        ▼
notification schema
        │
notification_logs
```
Each service only reads and writes its own schema.

## ER Diagram
> Refer: `docs/images/er-diagram.png`

## Architecture Decision Record (ADR)

### Why PostgreSQL instead of MongoDB?

```
CloudCart models a highly relational domain where users place orders, orders contain multiple order items, products belong to categories, and payments are associated with orders. These relationships benefit from foreign key constraints, SQL joins, and ACID transactions to maintain data consistency. PostgreSQL provides strong transactional guarantees and native referential integrity, making it a natural choice. While MongoDB supports references and multi-document transactions, it is better suited to domains with flexible or evolving document structures. For this e-commerce application, PostgreSQL aligns better with the data model and consistency requirements.
```

## Final Database Structure

```
Amazon RDS PostgreSQL
│
└── cloudcart database
    │
    ├── auth schema
    ├── product schema
    ├── orders schema
    └── notification schema

```

This is what we have chosen for CloudCart.
```
1 RDS Instance
    ↓
1 PostgreSQL Database
    ↓
4 Schemas
    ↓
Service-owned Tables
```

And each service gets a separate database user:

```
auth_service_user
    └── auth schema

product_service_user
    └── product schema

order_service_user
    └── orders schema

notification_service_user
    └── notification schema
```

Database with schemas with tables.

```
Amazon RDS PostgreSQL
│
└── cloudcart
    │
    ├── auth
    │   ├── users
    │   ├── roles
    │   └── refresh_tokens
    │
    ├── product
    │   ├── categories
    │   ├── products
    │   └── inventory
    │
    ├── orders
    │   ├── orders
    │   ├── order_items
    │   └── payments
    │
    └── notification
        └── notification_logs
```

### Important Microservices Rule

**One RDS PostgreSQL instance → one ```cloudcart``` database → multiple schemas → one database user per microservice with access only to its own schema.**