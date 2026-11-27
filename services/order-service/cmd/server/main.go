package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vidhyashekar/cloudcart/services/order-service/internal/app"
)

// main is the entry point of the order service application. It initializes the application and starts the HTTP server.
func main() {
	// Listen for Ctrl+C and Kubernetes termination signals.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	application, err := app.New()
	if err != nil {
		log.Printf("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	if err := application.Run(ctx); err != nil {
		log.Printf("Application stopped with error: %v", err)
		os.Exit(1)
	}
}
