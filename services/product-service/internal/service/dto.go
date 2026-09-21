package service

// CreateProductRequest represents the request payload for creating a new product.
type CreateProductRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	CategoryID    uint    `json:"category_id"`
}

// UpdateProductRequest represents the request payload for updating an existing product.
type UpdateProductRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	CategoryID    uint    `json:"category_id"`
}

// ProductResponse represents the response payload for a product.
type ProductResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	CategoryID    uint    `json:"category_id"`
}
