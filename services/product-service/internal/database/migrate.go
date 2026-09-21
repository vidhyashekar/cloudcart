package database

import (
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/model"
	"gorm.io/gorm"
)

// Migrate creates the necessary database schema and tables for the product service.
func Migrate(db *gorm.DB) error {
	if err := db.Exec(
		`CREATE SCHEMA IF NOT EXISTS product`,
	).Error; err != nil {
		return err
	}

	return db.AutoMigrate(
		&model.Category{},
		&model.Product{},
	)
}
