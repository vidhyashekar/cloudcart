package model

import "time"

type Product struct {
	ID            uint    `gorm:"primaryKey"`
	Name          string  `gorm:"size:255;not null"`
	Description   string  `gorm:"type:text"`
	Price         float64 `gorm:"not null"`
	StockQuantity int     `gorm:"not null"`

	CategoryID uint     `gorm:"not null"`
	Category   Category `gorm:"foreignKey:CategoryID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Product) TableName() string {
	return "product.products"
}
