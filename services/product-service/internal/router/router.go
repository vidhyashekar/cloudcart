package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/middleware"
)

func RegisterRoutes(productHandler *handler.ProductHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	productRoutes := router.Group("/api/v1/products")
	{
		productRoutes.GET("", productHandler.GetAll)
		productRoutes.GET("/:id", productHandler.GetByID)

		protected := productRoutes.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))

		protected.POST("", productHandler.Create)
		protected.PUT("/:id", productHandler.Update)
		protected.DELETE("/:id", productHandler.Delete)
	}

	return router
}
