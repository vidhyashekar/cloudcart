package repository

import (
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/model"
	"gorm.io/gorm"
)

// orderRepository is the concrete implementation of the OrderRepository interface.
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance of orderRepository with the provided database connection.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
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

// CreateWithTransaction creates an order along with its associated items in a single transaction.
func (r *orderRepository) CreateWithTransaction(order *model.Order, items []model.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil
	})
}
