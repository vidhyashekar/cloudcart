package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/middleware"
)

// RegisterRoutes sets up the routes for the order service, including health checks and order-related endpoints. It also applies authentication middleware to the order routes.
func RegisterRoutes(orderHandler *handler.OrderHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

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
