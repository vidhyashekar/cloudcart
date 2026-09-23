package app

import (
	"github.com/sirupsen/logrus"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/database"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	DB     *gorm.DB
}

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

	return &App{
		Config: cfg,
		DB:     db,
	}, nil
}
