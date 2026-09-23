package database

import (
	"fmt"

	"github.com/vidhyashekar/cloudcart/services/order-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect establishes a connection to the PostgreSQL database using the provided configuration.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
