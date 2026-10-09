package event

// OrderCreatedEvent represents the event structure for an order creation event.
type OrderCreatedEvent struct {
	EventType   string  `json:"event_type"` // Type of the event, e.g., "order.created", "order.confirmed", "order.canceled"
	OrderID     uint    `json:"order_id"`
	UserID      uint    `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
}
