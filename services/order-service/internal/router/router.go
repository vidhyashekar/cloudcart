package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/middleware"
)

// RegisterRoutes sets up health checks and order-related endpoints.
func RegisterRoutes(
	orderHandler *handler.OrderHandler,
	healthHandler *handler.HealthHandler,
	jwtSecret string,
) *gin.Engine {

	router := gin.Default()

	// Health check routes (no authentication required).
	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)

	// Protected order routes.
	orderRoutes := router.Group("/api/v1/orders")

	orderRoutes.Use(
		middleware.AuthMiddleware(jwtSecret),
	)

	{
		orderRoutes.POST("", orderHandler.Create)
		orderRoutes.GET("", orderHandler.GetAll)
		orderRoutes.GET("/:id", orderHandler.GetByID)
	}

	return router
}
