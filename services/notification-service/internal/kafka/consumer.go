package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/event"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/service"
)

// Consumer represents a Kafka consumer that listens for order events and processes them.
type Consumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	service  *service.NotificationService
}

// NewConsumer creates a new Kafka consumer for the specified topic and group ID.
func NewConsumer(brokers string, topic string, groupID string, notificationService *service.NotificationService) (*Consumer, error) {

	config := sarama.NewConfig()

	config.Consumer.Group.Rebalance.Strategy =
		sarama.NewBalanceStrategyRoundRobin()

	config.Consumer.Offsets.Initial =
		sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(
		[]string{brokers},
		groupID,
		config,
	)

	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
		service:  notificationService,
	}, nil
}

// consumerHandler is a custom implementation of sarama.ConsumerGroupHandler
type consumerHandler struct {
	service *service.NotificationService
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		log.Printf(
			"Kafka message received: topic=%s partition=%d offset=%d",
			message.Topic,
			message.Partition,
			message.Offset,
		)

		var orderEvent event.OrderCreatedEvent

		if err := json.Unmarshal(
			message.Value,
			&orderEvent,
		); err != nil {

			log.Printf(
				"Failed to deserialize Kafka message: %v",
				err,
			)

			continue
		}

		if err := h.service.HandleOrderCreated(
			orderEvent,
		); err != nil {

			log.Printf(
				"Failed to process order event: %v",
				err,
			)

			continue
		}

		session.MarkMessage(message, "")

		log.Printf(
			"Successfully processed order.created event for order %d",
			orderEvent.OrderID,
		)
	}

	return nil
}

// Start begins consuming messages from the Kafka topic and processes them using the consumerHandler.
func (c *Consumer) Start() {

	handler := &consumerHandler{
		service: c.service,
	}

	for {

		err := c.consumer.Consume(
			context.Background(),
			[]string{c.topic},
			handler,
		)

		if err != nil {
			log.Printf(
				"Kafka consumer error: %v",
				err,
			)
		}
	}
}
