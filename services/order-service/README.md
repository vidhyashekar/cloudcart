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

| Kafka term        | CloudCart                       |
| ----------------- | ------------------------------- |
| **Producer**      | Order Service                   |
| **Consumer**      | Notification Service            |
| **Topic**         | `order-events`                  |
| **Message/Event** | `order.created`                 |
| **Partition**     | 3 partitions in our topic       |
| **Broker**        | Our `cloudcart-kafka` container |


## Kafka

### Producer

A Producer sends messages/events to Kafka.

### Consumer

A Consumer reads messages from Kafka.

### Topic

A topic is like a named stream/category where events are stored.

For example:

```order-events```

We can publish:

```
order.created
order.cancelled
order.confirmed
```

to that topic.

Think of a topic approximately like:

a named channel for related events

### Kafka message

```
{
  "event_type": "order.created",
  "order_id": 101,
  "user_id": 25,
  "total_amount": 2099.98
}
```

Kafka doesn't really care about the business meaning of this JSON.

It essentially stores/transports the message.

### Partition
A topic can be divided into partitions.

For example:
```
order-events
│
├── Partition 0
├── Partition 1
└── Partition 2
```
Instead of having one giant stream, Kafka distributes messages across partitions.
This allows Kafka to handle large amounts of traffic and allows consumers to process partitions in parallel.

For our project, we created:
partitions = 3

### Kafka broker
A Kafka cluster might look like:
```
Kafka Cluster
│
├── Broker 1
├── Broker 2
└── Broker 3
```

Each broker can store and serve Kafka data.

For our local development:
```
Kafka Cluster
      │
      └── Broker 1
```

We're running only one Kafka broker.
That's why our Docker setup has:
```
KAFKA_NODE_ID: 1
```
