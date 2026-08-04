package main

import (
	"net/http"

	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/config"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/logger"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/router"
)

func main() {
	// Load configuration
	config := config.LoadConfig()

	// Initialize logger
	log := logger.New()

	// Register routes
	router.RegisterRoutes()

	log.Infof("Starting Auth-service on port %s", config.Port)

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal(err)
	}
}
