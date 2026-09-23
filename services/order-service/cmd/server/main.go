package main

import (
	"log"

	"github.com/vidhyashekar/cloudcart/services/order-service/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Order Service database initialized successfully")

	_ = application
}
