package outbox

import (
	"log"
	"time"

	"github.com/vidhyashekar/cloudcart/services/order-service/internal/kafka"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/repository"
)

// Publisher is responsible for publishing outbox events to Kafka.
type Publisher struct {
	Repository repository.OutboxRepository
	Producer   *kafka.Producer
}

// NewPublisher creates a new instance of Publisher.
func NewPublisher(repository repository.OutboxRepository, producer *kafka.Producer) *Publisher {
	return &Publisher{
		Repository: repository,
		Producer:   producer,
	}
}

// Start begins the process of publishing outbox events to Kafka at regular intervals.
func (p *Publisher) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Outbox publisher started")

	// Continuously listen for ticker events and publish pending outbox events.
	for range ticker.C {
		p.publishPendingEvents()
	}
}

// publishPendingEvents fetches pending outbox events from the repository and publishes them to Kafka.
func (p *Publisher) publishPendingEvents() {
	events, err := p.Repository.FindPendingEvents(10)

	if err != nil {
		log.Printf(
			"failed to fetch pending outbox events: %v",
			err,
		)
		return
	}

	if len(events) == 0 {
		log.Println("No pending outbox events")
		return
	}

	for _, outboxEvent := range events {

		err := p.Producer.Publish(
			outboxEvent.EventType,
			[]byte(outboxEvent.Payload),
		)

		if err != nil {
			log.Printf(
				"failed to publish outbox event %d: %v",
				outboxEvent.ID,
				err,
			)

			continue
		}

		if err := p.Repository.MarkAsPublished(
			outboxEvent.ID,
		); err != nil {
			log.Printf(
				"failed to mark outbox event %d as published: %v",
				outboxEvent.ID,
				err,
			)

			continue
		}

		log.Printf(
			"outbox event %d published successfully",
			outboxEvent.ID,
		)
	}
}
