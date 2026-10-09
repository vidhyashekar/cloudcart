package kafka

import (
	"github.com/IBM/sarama"
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

// Publish sends a message with the specified event type and payload to the Kafka topic. It constructs a ProducerMessage and sends it using the underlying Kafka producer. If any error occurs during sending, it returns the error.
func (p *Producer) Publish(eventType string, payload []byte) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(eventType),
		Value: sarama.ByteEncoder(payload),
	}

	_, _, err := p.producer.SendMessage(message)

	return err
}

// Close closes the Kafka producer.
func (p *Producer) Close() error {
	return p.producer.Close()
}
