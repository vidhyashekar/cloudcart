package app

import (
	"github.com/sirupsen/logrus"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/client"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/database"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/kafka"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/repository"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/router"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/service"
	"gorm.io/gorm"
)

// App struct holds the configuration, database connection, and services for the application.
type App struct {
	Config       *config.Config
	DB           *gorm.DB
	OrderService service.OrderService
	OrderHandler *handler.OrderHandler
}

// New initializes the application by loading the configuration, setting up the logger, connecting to the database, migrating the schema, and initializing the order service.
func New() (*App, error) {
	cfg := config.Load()

	log := logrus.New()

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to the database successfully")

	if err := database.Migrate(db); err != nil {
		return nil, err
	}
	log.Info("Schema and tables migrated successfully")

	orderRepository := repository.NewOrderRepository(db)

	productClient := client.NewProductClient(
		cfg.ProductServiceURL,
	)

	kafkaProducer, err := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		return nil, err
	}

	orderService := service.NewOrderService(orderRepository, productClient, kafkaProducer)

	orderHandler := handler.NewOrderHandler(orderService)

	return &App{
		Config:       cfg,
		DB:           db,
		OrderService: orderService,
		OrderHandler: orderHandler,
	}, nil
}

// Run starts the Gin router and listens for incoming HTTP requests on the specified port.
func (a *App) Run() error {
	r := router.RegisterRoutes(a.OrderHandler, a.Config.JWTSecret)

	return r.Run(":" + a.Config.Port)
}
