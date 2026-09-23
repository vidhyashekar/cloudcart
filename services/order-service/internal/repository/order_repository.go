package repository

import "github.com/vidhyashekar/cloudcart/services/order-service/internal/model"

// OrderRepository defines the interface for interacting with the order data in the database.
type OrderRepository interface {
	Create(order *model.Order) error
	FindByID(id uint) (*model.Order, error)
	FindByUserID(userID uint) ([]model.Order, error)
	Update(order *model.Order) error
	CreateWithTransaction(order *model.Order, items []model.OrderItem) error
}
