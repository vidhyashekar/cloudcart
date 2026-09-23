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

// ProductQuery represents the query parameters for filtering and paginating products.
type ProductQuery struct {
	Search     string
	CategoryID uint
	Page       int
	Limit      int
}

// ProductListResponse represents the response payload for a list of products with pagination information.
type ProductListResponse struct {
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"total_pages"`
	Products   []ProductResponse `json:"products"`
}

// UpdateStockRequest represents the request payload for updating the stock quantity of a product.
type UpdateStockRequest struct {
	Quantity int `json:"quantity"`
}
