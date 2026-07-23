# API Design Document

## API Boundaries

Our external request flow is:

```
Client
  │
  ▼
Application Load Balancer
  │
  ▼
Nginx Ingress Controller
  │
  ▼
CloudCart API Gateway (Go)
  │
  ├──► Auth Service
  ├──► Product Service
  ├──► Order Service
  └──► Notification Service
```

The client should not directly call individual microservices.

1. Auth Service APIs

Base path:
```/api/v1/auth```

| Method | Endpoint         | Purpose                     |
| ------ | ---------------- | --------------------------- |
| `POST` | `/register`      | Create a new user           |
| `POST` | `/login`         | Authenticate user           |
| `POST` | `/refresh-token` | Generate a new access token |
| `POST` | `/logout`        | Invalidate refresh token    |
| `GET`  | `/me`            | Get current user details    |

Example:

```http
POST /api/v1/auth/login
```

Request:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "access_token": "jwt-access-token",
  "refresh_token": "refresh-token",
  "expires_in": 3600
}
```

2. Product Service APIs

Base path:
```/api/v1```

| Method   | Endpoint         | Purpose             |
| -------- | ---------------- | ------------------- |
| `GET`    | `/products`      | List products       |
| `GET`    | `/products/{id}` | Get product details |
| `POST`   | `/products`      | Create product      |
| `PUT`    | `/products/{id}` | Update product      |
| `DELETE` | `/products/{id}` | Delete product      |
| `GET`    | `/categories`    | List categories     |

Example product search

Request:
```http
GET /api/v1/products?category=electronics&page=1&limit=20
```

Response:
```json
{
  "data": [
    {
      "id": 1,
      "name": "Laptop",
      "price": 75000,
      "category_id": 2,
      "stock": 10
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 1
}
```

3. Order Service APIs

Base path: ```/api/v1```

| Method  | Endpoint              | Purpose             |
| ------- | --------------------- | ------------------- |
| `POST`  | `/orders`             | Create order        |
| `GET`   | `/orders/{id}`        | Get order details   |
| `GET`   | `/orders`             | Get user's orders   |
| `PATCH` | `/orders/{id}/status` | Update order status |
| `POST`  | `/orders/{id}/cancel` | Cancel order        |


Create order
```
POST /api/v1
Authorization: Bearer <access-token>
```

Request:
```json
{
  "items": [
    {
      "product_id": 101,
      "quantity": 2
    }
  ]
}
```

Response:
```json
{
  "order_id": 5001,
  "status": "PENDING",
  "total_amount": 150000
}
```

4. Notification Service APIs
For our architecture, notification processing will primarily be event-driven through Kafka.

For example:
```
Order Service
     │
     │ publishes
     ▼
Kafka
     │
     │ consumes OrderCreated
     ▼
Notification Service
```

The notification service can have:
```GET /api/v1/notifications```
to allow users to view thier notification history.

## Important design decision

### Synchronous REST
Use when the caller needs an immediate response:
```
Client
  │
  ▼
Order Service
  │
  └── REST ──► Product Service
                │
                └── Check product and stock
```

### Asynchronous Kafka
Use for background processing:
```
Order Service
      │
      ▼
OrderCreated Event
      │
      ▼
Kafka
      │
      ▼
Notification Service
      │
      ▼
Send notification
```

So our order creation flow will be:
```
Client
  │
  ▼
API Gateway
  │
  ▼
Order Service
  │
  ├──► Validate JWT
  │
  ├──► Call Product Service
  │       └── Validate product and stock
  │
  ├──► Save Order
  │
  └──► Publish OrderCreated Event
              │
              ▼
           Kafka
              │
              ▼
      Notification Service
```

## Error Response Format

All APIs use a consistent error response format:

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Invalid request parameters",
    "details": []
  }
}
```
Example:

```json
{
  "error": {
    "code": "PRODUCT_NOT_FOUND",
    "message": "Product not found",
    "details": []
  }
}
```

## HTTP Status Codes

| Status Code | Meaning |
|-------------|---------|
| 200 | Successful request |
| 201 | Resource created |
| 204 | Successful request with no response body |
| 400 | Bad request |
| 401 | Unauthenticated |
| 403 | Forbidden |
| 404 | Resource not found |
| 409 | Conflict |
| 422 | Unprocessable entity |
| 500 | Internal server error |


## API Gateway Responsibilities

The CloudCart API Gateway is responsible for:

- Routing requests to appropriate services
- JWT validation
- Request authentication
- Request logging
- Rate limiting
- Consistent error handling
- Correlation ID propagation