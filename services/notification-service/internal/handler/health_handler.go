package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/notification-service/internal/kafka"
)

type HealthHandler struct {
	Consumer *kafka.Consumer
}

func NewHealthHandler(consumer *kafka.Consumer) *HealthHandler {
	return &HealthHandler{
		Consumer: consumer,
	}
}

func (h *HealthHandler) Live(c *gin.Context) {
	if !h.Consumer.IsRunning() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	if !h.Consumer.IsReady() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
			"checks": gin.H{
				"kafka_consumer": "DOWN",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
		"checks": gin.H{
			"kafka_consumer": "UP",
		},
	})
}
