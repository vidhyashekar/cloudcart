package model

import "time"

// OutboxEvent represents an event that is stored in the outbox table for eventual processing and publishing to Kafka.
type OutboxEvent struct {
	ID          uint   `gorm:"primaryKey"`
	EventType   string `gorm:"size:100;not null"` // Type of the event, e.g., "order.created"
	AggregateID uint   `gorm:"not null"`          // ID of the order associated with the event
	Payload     string `gorm:"type:jsonb;not null"`
	Status      string `gorm:"size:20;not null;default:PENDING"` // PENDING or PUBLISHED
	CreatedAt   time.Time
	PublishedAt *time.Time // Nullable, set when the event is published
}

// TableName specifies the table name for the OutboxEvent model.
func (OutboxEvent) TableName() string {
	return "order.outbox_events"
}
