package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/auth"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/middleware"
)

// RegisterRoutes sets up the routes for the authentication service, including health checks and authentication-related endpoints. It also applies authentication middleware to protected routes.
func RegisterRoutes(
	authHandler *handler.AuthHandler,
	healthHandler *handler.HealthHandler,
	jwtManager *auth.JWTManager,
) *gin.Engine {

	router := gin.Default()

	// Health check endpoints
	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)

	// Public authentication endpoints
	router.POST("/api/v1/auth/register", authHandler.Register)
	router.POST("/api/v1/auth/login", authHandler.Login)
	router.POST("/api/v1/auth/refresh", authHandler.Refresh)

	// Protected authentication endpoints
	authRoutes := router.Group("/api/v1/auth")
	authRoutes.Use(middleware.AuthMiddleware(jwtManager))
	{
		authRoutes.GET("/me", authHandler.Me)
	}

	return router
}
