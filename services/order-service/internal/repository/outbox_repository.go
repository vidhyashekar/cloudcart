package repository

import (
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/model"
	"gorm.io/gorm"
)

// OutboxRepository defines the interface for interacting with the outbox events in the database.
type OutboxRepository interface {
	Create(tx *gorm.DB, event *model.OutboxEvent) error
	FindPendingEvents(limit int) ([]model.OutboxEvent, error)
	MarkAsPublished(id uint) error
}
