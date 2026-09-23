package client

// ProductResponse represents the response structure for a product.
type ProductResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
}

// ProductClient defines the interface for interacting with the product service.
type ProductClient interface {
	GetProduct(id uint) (*ProductResponse, error)
}
