# CloudCart Services

We decided on four services.

``` CloudCart

├── Authentication Service
├── Product Service
├── Order Service
└── Notification Service

```

## Authentication Service

Responsible for identity.
 
### Owns
```
- Users
- Roles
- Passwords
- JWT Tokens
- Refresh Token
```

### APIs

```
POST /register
POST /login
POST /refresh
GET /profile
```

### Database
```
- users
- roles
- refresh_tokens
```

Produces no Kafka events.


## Product Service

Responsible for catalog.

### Owns
```
- Products
- Categories
- Inventory
- Search
```

### APIs

```
GET /products
GET /products/{id}
POST /products
PUT /products/{id}
DELETE /products/{id}
```

### Database
```
- products
- categories
- inventory
```

## Order Service

Responsible for:

```
- Cart
- Orders
- Payments
- Checkout
- Order Status
```


### APIs

```
POST /orders
GET /orders/{id}
GET /orders
PUT /orders/{id}/cancel
```

### Database
```
- orders
- order_items
- payments
```

### Publishes Kafka Events
```
OrderCreated
OrderCancelled
OrderCompleted
```

## Notification Service

Responsible for:

```
- Emails
- SMS
- Push Notifications 
```

### Consumes

```
OrderCreated
OrderCancelled
OrderCompleted
```

### Database

```
- notification_logs
```

## Communication Matrix

| From         | To      | Type          |
| ------------ | ------- | ------------- |
| API Gateway  | Auth    | REST          |
| API Gateway  | Product | REST          |
| API Gateway  | Order   | REST          |
| Order        | Auth    | REST          |
| Order        | Product | REST          |
| Order        | Kafka   | Publish Event |
| Notification | Kafka   | Consume Event |


## Service Ownership

| Service      | Owns          |
| ------------ | ------------- |
| Auth         | Users         |
| Product      | Products      |
| Order        | Orders        |
| Notification | Notifications |


## Design Principles

- Each service owns its own database.
- Services never access another service's database directly.
- Synchronous communication uses REST APIs.
- Asynchronous communication uses Apache Kafka.
- Each service can be independently deployed and scaled.