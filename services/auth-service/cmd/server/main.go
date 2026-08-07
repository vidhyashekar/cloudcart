package main

import (
	"log"

	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/app"
)

func main() {
	// Load dependencies and initialize the application
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
