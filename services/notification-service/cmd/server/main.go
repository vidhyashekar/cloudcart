package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/kafka"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/router"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/service"
)

func main() {
	cfg := config.Load()

	notificationService := service.NewNotificationService()

	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaGroupID,
		notificationService,
	)
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	healthHandler := handler.NewHealthHandler(consumer)

	r := router.RegisterRoutes(healthHandler)

	server := &http.Server{
		Addr:              ":8083",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	consumerDone := make(chan struct{})

	go func() {
		defer close(consumerDone)
		consumer.Start(ctx)
	}()

	serverDone := make(chan error, 1)

	go func() {
		log.Println("Notification health server listening on :8083")

		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverDone <- err
		} else {
			serverDone <- nil
		}
	}()

	log.Println("Notification Service started")
	log.Printf("Listening to Kafka topic: %s", cfg.KafkaTopic)

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err := <-serverDone:
		if err != nil {
			log.Printf("Health server error: %v", err)
		}
	}

	// Stop accepting new HTTP requests.
	log.Println("Shutting down health server...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Health server shutdown error: %v", err)
		_ = server.Close()
	}

	// Stop Kafka consumer.
	log.Println("Stopping Kafka consumer...")
	stop()

	<-consumerDone

	// Close Kafka consumer group.
	log.Println("Closing Kafka consumer...")

	if err := consumer.Close(); err != nil {
		log.Printf("Kafka consumer close error: %v", err)
	}

	log.Println("Notification Service shutdown completed")
}
