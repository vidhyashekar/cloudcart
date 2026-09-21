package model

import "time"

// Category represents a product category in the system.
type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null;unique"`
	Description string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName specifies the table name for the Category model.
func (Category) TableName() string {
	return "product.categories"
}
