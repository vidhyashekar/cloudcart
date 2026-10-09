│  Order Service  │
                    └────────┬────────┘
                             │
                             │ order.created
                             ▼
                    ┌─────────────────┐
                    │      Kafka      │
                    │  order-events   │
                    └────────┬────────┘
                             │
                             │ consume
                             ▼
                ┌─────────────────────────┐
                │  Notification Service  │
                └────────────┬────────────┘
                             │
                             ▼
                    Send Notification
                    (currently simulated)