package repository

import "github.com/vidhyashekar/cloudcart/services/product-service/internal/model"

type ProductFilter struct {
	Search     string
	CategoryID uint
	Page       int
	Limit      int
}

// ProductRepository defines the interface for interacting with product data in the database.
type ProductRepository interface {
	Create(product *model.Product) error
	FindAll(filter ProductFilter) ([]model.Product, int64, error)
	FindByID(id uint) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id uint) error
	GetCategoryByID(categoryID uint) (*model.Category, error)
	UpdateStock(id uint, quantity int) error
}
