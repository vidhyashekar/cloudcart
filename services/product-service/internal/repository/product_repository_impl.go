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
func (r *productRepository) FindAll(filter ProductFilter) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{})

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where(
			"name ILIKE ? OR description ILIKE ?",
			search,
			search,
		)
	}

	if filter.CategoryID != 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit

	err := query.
		Preload("Category").
		Offset(offset).
		Limit(filter.Limit).
		Order("id DESC").
		Find(&products).Error

	return products, total, err
}

// Update modifies an existing product in the database.
func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

// Delete removes a product from the database by its ID.
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

// GetCategoryByID retrieves a category by its ID from the database.
func (r *productRepository) GetCategoryByID(categoryID uint) (*model.Category, error) {
	var category model.Category

	if err := r.db.First(&category, categoryID).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

// UpdateStock updates the stock quantity of a product in the database. It ensures that the stock quantity does not go below zero.
func (r *productRepository) UpdateStock(id uint, quantity int) error {
	result := r.db.
		Model(&model.Product{}).
		Where("id = ? AND stock_quantity >= ?", id, quantity).
		UpdateColumn(
			"stock_quantity",
			gorm.Expr("stock_quantity - ?", quantity),
		)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
