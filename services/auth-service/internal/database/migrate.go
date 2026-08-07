package database

import (
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Create auth schema if it doesn't exist
	if err := db.Exec(`CREATE SCHEMA IF NOT EXISTS auth`).Error; err != nil {
		return err
	}

	// Create/Update tables
	if err := db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.RefreshToken{},
	); err != nil {
		return err
	}

	return nil
}
