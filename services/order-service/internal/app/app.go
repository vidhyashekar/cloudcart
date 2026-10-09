package app

import (
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
	}, nil
}

// Run starts the Gin router and listens for incoming HTTP requests on the specified port.
func (a *App) Run() error {
	// Start the outbox publisher in a separate goroutine to continuously process and publish events in the background.
	go a.OutboxPublisher.Start()

	// register routes
	r := router.RegisterRoutes(a.OrderHandler, a.Config.JWTSecret)

	// Start the HTTP server and listen on the configured port.
	return r.Run(":" + a.Config.Port)
}
