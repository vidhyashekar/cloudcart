package database

import (
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/model"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {

	var customerRole model.Role

	result := db.
		Where("name = ?", "CUSTOMER").
		First(&customerRole)

	if result.Error == nil {
		return nil
	}

	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	customerRole = model.Role{
		Name: "CUSTOMER",
	}

	return db.Create(&customerRole).Error
}
