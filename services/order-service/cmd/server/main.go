package main

import (
	"log"

	"github.com/vidhyashekar/cloudcart/services/order-service/internal/app"
)

// main is the entry point of the order service application. It initializes the application and starts the HTTP server.
func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
