package model

// ProcessPaymentRequest represents the request to process a payment
type ProcessPaymentRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// OrderCompleteEvent represents the event sent when an order is completed
type OrderCompleteEvent struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// CreatedOrderEvent represents the event sent when an order is created
type CreatedOrderEvent struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
}
