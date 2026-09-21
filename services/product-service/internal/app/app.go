package app

import (
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/database"
	"gorm.io/gorm"
)

// This file is part of the product-service application. It initializes the application, loads configuration, connects to the database, and runs migrations.
type App struct {
	Config *config.Config
	DB     *gorm.DB
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

	return &App{
		Config: cfg,
		DB:     db,
	}, nil
}
