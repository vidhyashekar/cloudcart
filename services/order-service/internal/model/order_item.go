package model

import "time"

type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	UnitPrice float64 `gorm:"not null"`
	Subtotal  float64 `gorm:"not null"`

	Order Order `gorm:"foreignKey:OrderID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (OrderItem) TableName() string {
	return "order.order_items"
}
