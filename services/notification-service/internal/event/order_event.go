package event

type OrderCreatedEvent struct {
	EventType   string  `json:"event_type"`
	OrderID     uint    `json:"order_id"`
	UserID      uint    `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
}
