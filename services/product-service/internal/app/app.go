package app

import (
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/database"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/repository"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/router"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/service"
	"gorm.io/gorm"
)

// App struct holds the configuration, database connection, and product handler for the application.
type App struct {
	Config         *config.Config
	DB             *gorm.DB
	ProductHandler *handler.ProductHandler
}

// New initializes a new instance of the App struct, loading configuration, connecting to the database, and running migrations.
func New() (*App, error) {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(db); err != nil {
		return nil, err
	}

	if err := database.SeedCategories(db); err != nil {
		return nil, err
	}

	productRepository := repository.NewProductRepository(db)

	productService := service.NewProductService(
		productRepository,
	)

	productHandler := handler.NewProductHandler(
		productService,
	)

	return &App{
		Config:         cfg,
		DB:             db,
		ProductHandler: productHandler,
	}, nil
}

// Run starts the HTTP server and registers the routes for the application.
func (a *App) Run() error {
	r := router.RegisterRoutes(a.ProductHandler, a.Config.JWTSecret)

	return r.Run(":" + a.Config.Port)
}
