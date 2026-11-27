package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync/atomic"

	"github.com/IBM/sarama"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/event"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/service"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	service  *service.NotificationService

	ready   atomic.Bool
	running atomic.Bool
}

func NewConsumer(
	brokers string,
	topic string,
	groupID string,
	notificationService *service.NotificationService,
) (*Consumer, error) {

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

func (c *Consumer) IsReady() bool {
	return c.ready.Load()
}

func (c *Consumer) IsRunning() bool {
	return c.running.Load()
}

type consumerHandler struct {
	service *service.NotificationService
	ready   *atomic.Bool
}

func (h *consumerHandler) Setup(
	sarama.ConsumerGroupSession,
) error {
	h.ready.Store(true)
	log.Println("Kafka consumer session established")
	return nil
}

func (h *consumerHandler) Cleanup(
	sarama.ConsumerGroupSession,
) error {
	h.ready.Store(false)
	log.Println("Kafka consumer session ended")
	return nil
}

func (h *consumerHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

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

func (c *Consumer) Start(ctx context.Context) {
	c.running.Store(true)
	defer c.running.Store(false)
	defer c.ready.Store(false)

	handler := &consumerHandler{
		service: c.service,
		ready:   &c.ready,
	}

	for {
		if ctx.Err() != nil {
			return
		}

		err := c.consumer.Consume(
			ctx,
			[]string{c.topic},
			handler,
		)

		if ctx.Err() != nil {
			return
		}

		if err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
