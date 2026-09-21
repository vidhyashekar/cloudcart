package main

import (
	"log"

	"github.com/vidhyashekar/cloudcart/services/product-service/internal/app"
)

func main() {
	log.Println("Starting the server...")
	// Load dependencies and initialize the application
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Product Service database initialized successfully")

	_ = app

}
