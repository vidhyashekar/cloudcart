package service

import (
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/model"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/repository"
)

// ProductService defines the interface for product-related operations.
type ProductService interface {
	CreateProduct(request CreateProductRequest) (*ProductResponse, error)
	GetProducts() ([]ProductResponse, error)
	GetProductByID(id uint) (*ProductResponse, error)
	UpdateProduct(id uint, request UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(id uint) error
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
func (s *productService) GetProducts() ([]ProductResponse, error) {
	products, err := s.productRepository.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]ProductResponse, 0, len(products))

	for _, product := range products {
		responses = append(responses, *toProductResponse(&product))
	}

	return responses, nil
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
