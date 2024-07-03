package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"order/model"
	"order/service"
)

// Handler is the main handler for the create order API
type Handler struct {
	orderService   service.OrderService
	paymentService service.PaymentService
}

func NewHandler(orderService service.OrderService, paymentService service.PaymentService) *Handler {
	return &Handler{
		orderService:   orderService,
		paymentService: paymentService,
	}
}

// HandleRequest handles both API Gateway and SQS events
func (h *Handler) HandleRequest(ctx context.Context, event json.RawMessage) (interface{}, error) {
	var apiEvent events.APIGatewayProxyRequest
	var sqsEvent events.SQSEvent

	// Try to unmarshal into APIGatewayProxyRequest
	if err := json.Unmarshal(event, &apiEvent); err == nil && apiEvent.HTTPMethod != "" {
		return h.CreateOrder(ctx, apiEvent)
	}

	// Try to unmarshal into SQSEvent
	if err := json.Unmarshal(event, &sqsEvent); err == nil && len(sqsEvent.Records) > 0 {
		return nil, h.CompleteOrder(ctx, sqsEvent)
	}

	// Unsupported event type
	log.Println("Unsupported event type")
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       "Unsupported event type",
	}, nil
}

// CreateOrder handles the create order API request
func (h *Handler) CreateOrder(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var order model.CreateOrderRequest
	err := json.Unmarshal([]byte(request.Body), &order)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	orderID, err := h.orderService.CreateOrder(ctx, order)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	event := model.CreateOrderEvent{
		OrderID:    orderID,
		TotalPrice: order.TotalPrice,
	}

	err = h.paymentService.NotifyPayment(ctx, event)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	order.ID = orderID
	responseBody, err := json.Marshal(order)
	if err != nil {
		log.Printf("Failed to marshal response body: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to marshal response body"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(responseBody)}, nil
}

// CompleteOrder handles the order completion event from SQS
func (h *Handler) CompleteOrder(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		var event model.OrderCompletedEvent
		err := json.Unmarshal([]byte(message.Body), &event)
		if err != nil {
			return err
		}

		err = h.orderService.CompleteOrder(ctx, event.OrderID, event.Status)
		if err != nil {
			return err
		}
	}

	return nil
}
