package app

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/database"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/router"
	"gorm.io/gorm"
)

// App struct holds the configuration, logger, and database connection for the application.
type App struct {
	Config *config.Config
	Logger *logrus.Logger
	DB     *gorm.DB
}

func New() (*App, error) {
	config := config.Load()

	log := logrus.New()

	db, err := database.Connect(config)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to the database successfully")

	err = database.Migrate(db)
	if err != nil {
		return nil, err
	}
	log.Info("Schema and tables migrated successfully")

	return &App{
		Config: config,
		Logger: log,
		DB:     db,
	}, nil
}

func (a *App) Run() error {
	router.RegisterRoutes()

	a.Logger.Infof("Auth Service started on port %s", a.Config.Port)
	return http.ListenAndServe(":"+a.Config.Port, nil)
}
