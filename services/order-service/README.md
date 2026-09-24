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
A → Atomicity → All or nothing 
C → Consistency → Valid state before/after 
I → Isolation → Concurrent transactions don't improperly interfere 
D → Durability → Committed data survives failures

Problem:

But there's one important consistency issue

We shouldn't simply do this:

1. Reduce stock
2. Create order

because if:

Reduce stock       ✅
Create order       ❌

we've reduced stock for an order that doesn't exist.

That's exactly the kind of distributed consistency problem we were discussing with transactions.

For our project, we'll handle this using a compensating action:

Get products
   ↓
Reduce stock
   ↓
Create order transaction
   │
   ├── success → publish order.created
   │
   └── failure → restore stock

We'll add the restore-stock operation next, and then integrate Kafka.

That will give us a practical introduction to the Saga/compensation pattern without overcomplicating the project.
```

The complete consistency flow
                 Order Service
                      │
              Get product details
                      │
                      ▼
              Check price/stock
                      │
                      ▼
              Decrease stock
                      │
             ┌────────┴────────┐
             │                 │
          Success            Failure
             │                 │
             ▼                 ▼
       Create Order        Restore stock
       Transaction
             │
        ┌────┴────┐
        │         │
     Success    Failure
        │         │
        ▼         ▼
     Continue   Restore
        │
        ▼
      Kafka


Kafka architecture
```
                     ┌───────────────────┐
                     │   Order Service   │
                     │                   │
                     │     Producer      │
                     └─────────┬─────────┘
                               │
                               │ order.created
                               ▼
                     ┌───────────────────┐
                     │      Kafka        │
                     │                   │
                     │   order-events    │
                     │                   │
                     │ ┌───┬───┬───┐     │
                     │ │ P0│ P1│ P2│     │
                     │ └───┴───┴───┘     │
                     └─────────┬─────────┘
                               │
                               ▼
                     ┌───────────────────┐
                     │ Notification      │
                     │ Service            │
                     │                   │
                     │     Consumer      │
                     └───────────────────┘

```
