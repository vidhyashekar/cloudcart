package database

import (
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/model"
	"gorm.io/gorm"
)

// SeedCategories checks if the "Electronics" category exists in the database. If it doesn't, it creates the category with a description of "Electronic products".
func SeedCategories(db *gorm.DB) error {
	var count int64

	if err := db.Model(&model.Category{}).
		Where("name = ?", "Electronics").
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return db.Create(&model.Category{
		Name:        "Electronics",
		Description: "Electronic products",
	}).Error
}
