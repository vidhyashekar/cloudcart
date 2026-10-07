package main

import (
	"log"

	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/kafka"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/service"
)

func main() {

	cfg := config.Load()

	notificationService :=
		service.NewNotificationService()

	// Create a Kafka consumer to listen for order events
	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaGroupID,
		notificationService,
	)

	if err != nil {
		log.Fatalf(
			"failed to create Kafka consumer: %v",
			err,
		)
	}

	log.Println(
		"Notification Service started",
	)

	log.Printf(
		"Listening to Kafka topic: %s",
		cfg.KafkaTopic,
	)

	//
	consumer.Start()
}
