package model

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	ID         string `json:"id,omitempty"`
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
}

// CreateOrderEvent represents the event sent when an order is created
type CreateOrderEvent struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
}

// OrderCompletedEvent represents the event sent when an order is completed
type OrderCompletedEvent struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
