package service

import (
	"errors"
	"fmt"

	"github.com/vidhyashekar/cloudcart/services/product-service/internal/model"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/repository"
	"gorm.io/gorm"
)

// ProductService defines the interface for product-related operations.
type ProductService interface {
	CreateProduct(request CreateProductRequest) (*ProductResponse, error)
	GetProducts(query ProductQuery) (*ProductListResponse, error)
	GetProductByID(id uint) (*ProductResponse, error)
	UpdateProduct(id uint, request UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(id uint) error
	UpdateStock(id uint, quantity int) error
}

// CreateProductRequest represents the request payload for creating a new product.
type productService struct {
	productRepository repository.ProductRepository
}

// NewProductService creates a new instance of productService with the provided ProductRepository.
func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}

// CreateProduct creates a new product in the database.
func (s *productService) CreateProduct(request CreateProductRequest) (*ProductResponse, error) {
	product := &model.Product{
		Name:          request.Name,
		Description:   request.Description,
		Price:         request.Price,
		StockQuantity: request.StockQuantity,
		CategoryID:    request.CategoryID,
	}

	if err := s.productRepository.Create(product); err != nil {
		return nil, err
	}

	return toProductResponse(product), nil
}

// GetProducts retrieves all products from the database.
func (s *productService) GetProducts(query ProductQuery) (*ProductListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.Limit < 1 {
		query.Limit = 10
	}

	if query.Limit > 100 {
		query.Limit = 100
	}

	filter := repository.ProductFilter{
		Search:     query.Search,
		CategoryID: query.CategoryID,
		Page:       query.Page,
		Limit:      query.Limit,
	}

	products, total, err := s.productRepository.FindAll(filter)
	if err != nil {
		return nil, err
	}

	responses := make([]ProductResponse, 0, len(products))

	for _, product := range products {
		responses = append(responses, *toProductResponse(&product))
	}

	totalPages := int((total + int64(query.Limit) - 1) / int64(query.Limit))

	return &ProductListResponse{
		Products:   responses,
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// GetProductByID retrieves a product by its ID.
func (s *productService) GetProductByID(id uint) (*ProductResponse, error) {
	product, err := s.productRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return toProductResponse(product), nil
}

// UpdateProduct updates an existing product by its ID.
func (s *productService) UpdateProduct(id uint, request UpdateProductRequest) (*ProductResponse, error) {
	product, err := s.productRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	_, err = s.productRepository.GetCategoryByID(request.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category with ID %d not found", request.CategoryID)
		}
		return nil, err
	}

	product.Name = request.Name
	product.Description = request.Description
	product.Price = request.Price
	product.StockQuantity = request.StockQuantity
	product.CategoryID = request.CategoryID

	if err := s.productRepository.Update(product); err != nil {
		return nil, err
	}

	return toProductResponse(product), nil
}

// DeleteProduct removes a product by its ID.
func (s *productService) DeleteProduct(id uint) error {
	_, err := s.productRepository.FindByID(id)
	if err != nil {
		return fmt.Errorf("product with ID %d not found", id)
	}
	return s.productRepository.Delete(id)
}

// toProductResponse converts a Product model to a ProductResponse DTO.
func toProductResponse(product *model.Product) *ProductResponse {
	return &ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		CategoryID:    product.CategoryID,
	}
}

// UpdateStock updates the stock quantity of a product by its ID.
func (s *productService) UpdateStock(id uint, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	return s.productRepository.UpdateStock(id, quantity)
}
