package database

import (
	"github.com/vidhyashekar/cloudcart/services/order-service/internal/model"
	"gorm.io/gorm"
)

// Migrate creates the necessary database schema and tables for the order service.
func Migrate(db *gorm.DB) error {
	if err := db.Exec(
		`CREATE SCHEMA IF NOT EXISTS "order"`,
	).Error; err != nil {
		return err
	}

	return db.AutoMigrate(
		&model.Order{},
		&model.OrderItem{},
	)
}
