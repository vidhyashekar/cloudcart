package repository

import (
	"encoding/json"

	"github.com/vidhyashekar/cloudcart/services/order-service/internal/event"
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/model"
	"gorm.io/gorm"
)

// orderRepository is the concrete implementation of the OrderRepository interface.
type orderRepository struct {
	db               *gorm.DB
	outboxRepository OutboxRepository
}

// NewOrderRepository creates a new instance of orderRepository with the provided database connection.
func NewOrderRepository(db *gorm.DB, outboxRepository OutboxRepository) OrderRepository {
	return &orderRepository{
		db:               db,
		outboxRepository: outboxRepository,
	}
}

// Create inserts a new order into the database.
func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

// FindByID retrieves an order by its ID, including its associated items.
func (r *orderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order

	err := r.db.
		Preload("Items").
		First(&order, id).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// FindByUserID retrieves all orders associated with a specific user ID.
func (r *orderRepository) FindByUserID(userID uint) ([]model.Order, error) {
	var orders []model.Order

	err := r.db.
		Preload("Items").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&orders).Error

	return orders, err
}

// Update updates an existing order in the database.
func (r *orderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

// CreateOrderWithOutbox creates an order, its associated items, and an outbox event within a single database transaction. If any operation fails, the entire transaction is rolled back.
func (r *orderRepository) CreateOrderWithOutbox(order *model.Order, items []model.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		// 1. Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Set generated Order ID on every item
		for i := range items {
			items[i].OrderID = order.ID
		}

		// 3. Create order items
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		// 4. Create event payload AFTER order ID exists
		eventPayload := event.OrderCreatedEvent{
			EventType:   "order.created",
			OrderID:     order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
		}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		// 5. Create outbox event
		outboxEvent := &model.OutboxEvent{
			EventType:   "order.created",
			AggregateID: order.ID,
			Payload:     string(payload),
			Status:      "PENDING",
		}

		if err := r.outboxRepository.Create(
			tx,
			outboxEvent,
		); err != nil {
			return err
		}

		return nil
	})
}
