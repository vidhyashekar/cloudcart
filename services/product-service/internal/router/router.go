package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/middleware"
)

func RegisterRoutes(
	productHandler *handler.ProductHandler,
	healthHandler *handler.HealthHandler,
	jwtSecret string,
) *gin.Engine {

	router := gin.Default()

	// Health check endpoints
	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)

	productRoutes := router.Group("/api/v1/products")
	{
		productRoutes.GET("", productHandler.GetAll)
		productRoutes.GET("/:id", productHandler.GetByID)

		protected := productRoutes.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))

		protected.POST("", productHandler.Create)
		protected.PUT("/:id", productHandler.Update)
		protected.DELETE("/:id", productHandler.Delete)
		protected.PATCH("/:id/stock/decrease", productHandler.DecreaseStock)
		protected.PATCH("/:id/stock/increase", productHandler.IncreaseStock)
	}

	return router
}
