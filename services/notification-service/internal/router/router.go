package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/handler"
)

func RegisterRoutes(
	healthHandler *handler.HealthHandler,
) *gin.Engine {

	router := gin.Default()

	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)

	return router
}
