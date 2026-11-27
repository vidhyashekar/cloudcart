package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	DB *gorm.DB
}

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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
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
