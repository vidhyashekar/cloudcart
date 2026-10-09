package repository

import (
	"time"

	"github.com/vidhyashekar/cloudcart/services/order-service/internal/model"
	"gorm.io/gorm"
)

// outboxRepository is the concrete implementation of the OutboxRepository interface.
type outboxRepository struct {
	db *gorm.DB
}

// NewOutboxRepository creates a new instance of outboxRepository with the provided GORM database connection.
func NewOutboxRepository(db *gorm.DB) OutboxRepository {
	return &outboxRepository{
		db: db,
	}
}

// Create inserts a new outbox event into the database within the provided transaction.
func (r *outboxRepository) Create(tx *gorm.DB, event *model.OutboxEvent) error {
	return tx.Create(event).Error
}

// FindPendingEvents retrieves a list of pending outbox events from the database, limited by the specified number.
func (r *outboxRepository) FindPendingEvents(limit int) ([]model.OutboxEvent, error) {
	var events []model.OutboxEvent

	err := r.db.
		Where("status = ?", "PENDING").
		Order("id ASC").
		Limit(limit).
		Find(&events).Error

	return events, err
}

// MarkAsPublished updates the status of an outbox event to "PUBLISHED" and sets the published_at timestamp.
func (r *outboxRepository) MarkAsPublished(id uint) error {
	now := time.Now()

	return r.db.
		Model(&model.OutboxEvent{}).
		Where("id = ? AND status = ?", id, "PENDING").
		Updates(map[string]interface{}{
			"status":       "PUBLISHED",
			"published_at": &now,
		}).Error
}
