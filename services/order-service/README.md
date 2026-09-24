# Order Service Flow

```
POST /orders
     ↓
Order Handler
     ↓
Order Service
     ↓
Product Client
     ↓
Product Service
     ├── validate product
     ├── get price
     └── check stock
     ↓
DB Transaction
     ├── create order
     ├── create order items
     └── reduce stock
     ↓
Kafka: order.created

```
