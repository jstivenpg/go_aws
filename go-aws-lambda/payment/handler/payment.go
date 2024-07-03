package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"payment/model"
	"payment/service"
)

// Handler is the main handler for the process payment API
type Handler struct {
	paymentService service.PaymentService
	orderService   service.OrderService
}

func NewHandler(paymentService service.PaymentService, orderService service.OrderService) *Handler {
	return &Handler{
		paymentService: paymentService,
		orderService:   orderService,
	}
}

// HandleRequest handles both API Gateway and SQS events
func (h *Handler) HandleRequest(ctx context.Context, event json.RawMessage) (interface{}, error) {
	var apiEvent events.APIGatewayProxyRequest
	var sqsEvent events.SQSEvent

	// Try to unmarshal into APIGatewayProxyRequest
	if err := json.Unmarshal(event, &apiEvent); err == nil && apiEvent.HTTPMethod != "" {
		return h.ProcessPayment(ctx, apiEvent)
	}

	// Try to unmarshal into SQSEvent
	if err := json.Unmarshal(event, &sqsEvent); err == nil && len(sqsEvent.Records) > 0 {
		return nil, h.CreatePayment(ctx, sqsEvent)
	}

	// Unsupported event type
	log.Println("Unsupported event type")
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       "Unsupported event type",
	}, nil
}

// ProcessPayment handles the process payment API request
func (h *Handler) ProcessPayment(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var payment model.ProcessPaymentRequest
	err := json.Unmarshal([]byte(request.Body), &payment)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	err = h.paymentService.ProcessPayment(ctx, payment)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	event := model.OrderCompleteEvent{
		OrderID: payment.OrderID,
		Status:  payment.Status,
	}

	err = h.orderService.NotifyOrderCompletion(ctx, event)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: request.Body}, nil
}

// CreatePayment handles the payment completion event from SQS
func (h *Handler) CreatePayment(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		var event model.CreatedOrderEvent
		err := json.Unmarshal([]byte(message.Body), &event)
		if err != nil {
			return err
		}

		err = h.paymentService.CreatePayment(ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}
