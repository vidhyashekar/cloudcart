package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/event"
)

// Producer is a struct that encapsulates the Kafka producer and the topic to which messages will be sent.
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer creates a new Kafka producer with the specified brokers and topic. It returns a pointer to the Producer and an error if any occurs during the creation of the producer.
func NewProducer(brokers string, topic string) (*Producer, error) {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(
		[]string{brokers},
		config,
	)

	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

// PublishOrderCreated publishes an OrderCreatedEvent to the Kafka topic. It serializes the event to JSON and sends it as a message to the specified topic. If any error occurs during serialization or sending, it returns the error.
func (p *Producer) PublishOrderCreated(event event.OrderCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(data),
		Key: sarama.StringEncoder(
			fmt.Sprintf("%d", event.OrderID),
		),
	}

	_, _, err = p.producer.SendMessage(message)

	return err
}

// Close closes the Kafka producer.
func (p *Producer) Close() error {
	return p.producer.Close()
}
