package repository

import "github.com/vidhyashekar/cloudcart/services/product-service/internal/model"

// ProductRepository defines the interface for interacting with product data in the database.
type ProductRepository interface {
	Create(product *model.Product) error
	FindAll() ([]model.Product, error)
	FindByID(id uint) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id uint) error
}
