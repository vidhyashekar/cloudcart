package main

import (
	"log"

	"github.com/vidhyashekar/cloudcart/services/product-service/internal/app"
)

// main is the entry point of the product-service application. It initializes the application and starts the HTTP server.
func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
