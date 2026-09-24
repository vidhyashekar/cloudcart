# Kafka local setup

## Step 1: Create Kafka infrastructure
```
infrastructure/
└── kafka/
    └── docker-compose.yml
```

## Step 2: Start Kafka
```
cd infrastructure/kafka
docker compose up -d
docker ps
docker logs cloudcart-kafka
```

## Step 3: Create Topic
Create:
```
docker exec cloudcart-kafka \
  /opt/kafka/bin/kafka-topics.sh \
  --create \
  --topic order-events \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1
  ```

  Verify:
  ```
  docker exec cloudcart-kafka \
  /opt/kafka/bin/kafka-topics.sh \
  --list \
  --bootstrap-server localhost:9092
  ```

  You should see: ```order-events```

### Why 3 partitions?

For learning, this lets us demonstrate that Kafka topics can be partitioned:

order-events
├── partition 0
├── partition 1
└── partition 2

Later, Kubernetes/AWS deployment can scale consumers based on partitions.

## Step 4: Test Kafka itself
Open terminal 1:
```
docker exec -it cloudcart-kafka \
  /opt/kafka/bin/kafka-console-consumer.sh \
  --topic order-events \
  --bootstrap-server localhost:9092 \
  --from-beginning
```
Keep it running.

Open terminal 2:
```
docker exec -it cloudcart-kafka \
  /opt/kafka/bin/kafka-console-producer.sh \
  --topic order-events \
  --bootstrap-server localhost:9092
```

Type: ```hello-cloudcart```

Terminal 1 should receive: ```hello-cloudcart```
That confirms Kafka itself is working.

We'll publish something like:
```
{
  "event_type": "order.created",
  "order_id": 101,
  "user_id": 25,
  "total_amount": 2099.98
}
```
to 
```order-events``` topic.

then:
```
Kafka
  │
  ▼
Notification Service
  │
  ▼
order.created
  │
  ▼
"Notification sent for Order #101"
```