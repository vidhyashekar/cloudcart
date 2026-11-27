package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles application liveness and readiness checks.
type HealthHandler struct {
	DB *gorm.DB
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		DB: db,
	}
}

// Live checks whether the application is alive.
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

// Ready checks whether the application is ready to accept traffic.
// It verifies database connectivity.
func (h *HealthHandler) Ready(c *gin.Context) {

	sqlDB, err := h.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
			"checks": gin.H{
				"database": "DOWN",
			},
		})
		return
	}

	// Set a timeout so the health check does not wait indefinitely.
	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		2*time.Second,
	)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
			"checks": gin.H{
				"database": "DOWN",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
		"checks": gin.H{
			"database": "UP",
		},
	})
}
