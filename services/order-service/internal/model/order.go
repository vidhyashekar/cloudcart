package model

import "time"

// Order represents a customer's order in the system.
type Order struct {
	ID          uint    `gorm:"primaryKey"`
	UserID      uint    `gorm:"not null"`
	Status      string  `gorm:"size:30;not null"`
	TotalAmount float64 `gorm:"not null"`

	Items []OrderItem `gorm:"foreignKey:OrderID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Order) TableName() string {
	return "order.orders"
}
