package service

// CreateOrderRequest represents the request payload for creating a new order.
type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}

// CreateOrderItemRequest represents the request payload for creating a new order item.
type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// OrderResponse represents the response payload for an order, including its items.
type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	Items       []OrderItemResponse `json:"items"`
}

// OrderItemResponse represents the response payload for an order item.
type OrderItemResponse struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}
