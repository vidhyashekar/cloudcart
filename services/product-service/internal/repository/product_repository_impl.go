package repository

import (
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/model"
	"gorm.io/gorm"
)

// productRepository is the concrete implementation of the ProductRepository interface.
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of productRepository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create adds a new product to the database.
func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

// FindByID retrieves a product by its ID from the database.
func (r *productRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product

	err := r.db.
		Preload("Category").
		First(&product, id).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// FindAll retrieves all products from the database.
func (r *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product

	err := r.db.
		Preload("Category").
		Find(&products).Error

	return products, err
}

// Update modifies an existing product in the database.
func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

// Delete removes a product from the database by its ID.
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}
