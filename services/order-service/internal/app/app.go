package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/client"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/database"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/kafka"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/outbox"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/repository"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/router"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/service"
	"gorm.io/gorm"
)

// App struct holds the configuration, database connection, and services for the application.
type App struct {
	Config          *config.Config
	DB              *gorm.DB
	OrderService    service.OrderService
	OrderHandler    *handler.OrderHandler
	OutboxPublisher *outbox.Publisher
	KafkaProducer   *kafka.Producer
}

// New initializes the application and wires all dependencies.
func New() (*App, error) {

	// 1. Load configuration.
	cfg := config.Load()

	// 2. Initialize logger.
	log := logrus.New()

	// 3. Connect to PostgreSQL.
	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}

	log.Info("Connected to the database successfully")

	// 4. Run database migrations.
	if err := database.Migrate(db); err != nil {
		return nil, err
	}

	log.Info("Schema and tables migrated successfully")

	// 5. Create Outbox Repository.
	outboxRepository := repository.NewOutboxRepository(db)

	// 6. Create Order Repository.
	// Order Repository owns the DB transaction that creates:
	// Order + Order Items + Outbox Event.
	orderRepository := repository.NewOrderRepository(
		db,
		outboxRepository,
	)

	// 7. Create Kafka Producer.
	// The producer is used by the Outbox Publisher,
	// NOT directly by Order Service.
	kafkaProducer, err := kafka.NewProducer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
	)
	if err != nil {
		return nil, err
	}

	// 8. Create Outbox Publisher.
	// It periodically reads PENDING events from the database
	// and publishes them to Kafka.
	outboxPublisher := outbox.NewPublisher(
		outboxRepository,
		kafkaProducer,
	)

	// 9. Create Product Service client.
	productClient := client.NewProductClient(
		cfg.ProductServiceURL,
	)

	// 10. Create Order Service.
	orderService := service.NewOrderService(
		orderRepository,
		productClient,
	)

	// 11. Create Order Handler.
	orderHandler := handler.NewOrderHandler(
		orderService,
	)

	// 12. Return fully initialized application.
	return &App{
		Config:          cfg,
		DB:              db,
		OrderService:    orderService,
		OrderHandler:    orderHandler,
		OutboxPublisher: outboxPublisher,
		KafkaProducer:   kafkaProducer,
	}, nil
}

// Run starts the HTTP server and background Outbox publisher.
func (a *App) Run(ctx context.Context) error {

	healthHandler := handler.NewHealthHandler(a.DB)

	// Register routes.
	r := router.RegisterRoutes(
		a.OrderHandler,
		healthHandler,
		a.Config.JWTSecret,
	)

	// Create HTTP server.
	server := &http.Server{
		Addr:              ":" + a.Config.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Create context for Outbox publisher.
	publisherCtx, cancelPublisher := context.WithCancel(context.Background())
	publisherDone := make(chan struct{})

	// Start Outbox publisher in a separate goroutine.
	go func() {
		defer close(publisherDone)
		a.OutboxPublisher.Start(publisherCtx)
	}()

	// Start HTTP server in a separate goroutine.
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("Order Service starting on port %s", a.Config.Port)

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	var runErr error

	// Wait for shutdown signal or HTTP server failure.
	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")

	case err := <-serverErr:
		runErr = fmt.Errorf("HTTP server failed: %w", err)
	}

	// 1. Gracefully shut down HTTP server.
	log.Println("Shutting down HTTP server...")

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown failed: %v", err)
		_ = server.Close()

		if runErr == nil {
			runErr = fmt.Errorf("HTTP server shutdown: %w", err)
		}
	}

	cancelShutdown()

	// 2. Stop Outbox publisher and wait for it to finish.
	log.Println("Stopping Outbox publisher...")
	cancelPublisher()
	<-publisherDone

	// 3. Close Kafka producer.
	log.Println("Closing Kafka producer...")
	if err := a.KafkaProducer.Close(); err != nil {
		log.Printf("Failed to close Kafka producer: %v", err)

		if runErr == nil {
			runErr = fmt.Errorf("close Kafka producer: %w", err)
		}
	}

	// 4. Close PostgreSQL connection.
	log.Println("Closing PostgreSQL connection...")

	sqlDB, err := a.DB.DB()
	if err != nil {
		if runErr == nil {
			runErr = fmt.Errorf("get SQL database connection: %w", err)
		}
	} else if err := sqlDB.Close(); err != nil {
		if runErr == nil {
			runErr = fmt.Errorf("close PostgreSQL connection: %w", err)
		}
	}

	log.Println("Order Service shutdown completed")

	return runErr
}
