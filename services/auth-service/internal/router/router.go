package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/auth"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/handler"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/middleware"
)

func RegisterRoutes(authHandler *handler.AuthHandler, jwtManager *auth.JWTManager) *gin.Engine {

	router := gin.Default()

	router.GET("/health", handler.HealthHandler)

	router.POST(
		"/api/v1/auth/register",
		authHandler.Register,
	)

	router.POST(
		"/api/v1/auth/login",
		authHandler.Login,
	)

	router.POST(
		"/api/v1/auth/refresh",
		authHandler.Refresh,
	)

	authRoutes := router.Group("/api/v1/auth")
	authRoutes.Use(middleware.AuthMiddleware(jwtManager))
	{
		authRoutes.GET("/me", authHandler.Me)
	}

	return router
}
