package router

import (
	"net/http"

	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/handler"
)

func RegisterRoutes() {
	http.HandleFunc("/health", handler.HealthHandler)
}
