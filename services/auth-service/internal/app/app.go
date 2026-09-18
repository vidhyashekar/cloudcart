package app

import (
	"github.com/sirupsen/logrus"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/auth"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/database"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/repository"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/router"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/service"
	"gorm.io/gorm"
)

// App struct holds the configuration, logger, and database connection for the application.
type App struct {
	Config         *config.Config
	Logger         *logrus.Logger
	DB             *gorm.DB
	UserRepository repository.UserRepository
	JWTManager     *auth.JWTManager
	AuthService    service.AuthService
	AuthHandler    *handler.AuthHandler
}

func New() (*App, error) {
	// Load configuration
	config := config.Load()

	// Logger setup
	log := logrus.New()

	// Database connection
	db, err := database.Connect(config)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to the database successfully")

	// Migrate schema and tables
	err = database.Migrate(db)
	if err != nil {
		return nil, err
	}
	log.Info("Schema and tables migrated successfully")

	// Seed roles
	err = database.SeedRoles(db)
	if err != nil {
		return nil, err
	}

	userRepository := repository.NewUserRepository(db)

	jwtManager := auth.NewJWTManager(
		config.JWTSecret,
		config.JWTExpiration,
	)

	authService := service.NewAuthService(userRepository, jwtManager)

	authHandler := handler.NewAuthHandler(authService)

	return &App{
		Config:         config,
		Logger:         log,
		DB:             db,
		UserRepository: userRepository,
		JWTManager:     jwtManager,
		AuthService:    authService,
		AuthHandler:    authHandler,
	}, nil
}

func (a *App) Run() error {
	router := router.RegisterRoutes(a.AuthHandler, a.JWTManager)

	a.Logger.Infof("Auth Service started on port %s", a.Config.Port)
	return router.Run(":" + a.Config.Port)
}
