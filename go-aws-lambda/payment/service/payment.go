package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"payment/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// PaymentService interface to handle payment operations
type PaymentService interface {
	// ProcessPayment processes a payment, updates the payment status and sends a message to the queue
	ProcessPayment(ctx context.Context, payment model.ProcessPaymentRequest) error
	// CreatePayment creates a payment record
	CreatePayment(ctx context.Context, event model.CreatedOrderEvent) error
}

type paymentService struct {
	dynamoDB *dynamodb.DynamoDB
	table    string
}

// NewPaymentService creates a new PaymentService
func NewPaymentService(dynamoDB *dynamodb.DynamoDB, table string) PaymentService {
	return &paymentService{
		dynamoDB: dynamoDB,
		table:    table,
	}
}

func (s *paymentService) ProcessPayment(ctx context.Context, payment model.ProcessPaymentRequest) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("PAYMENT#%s", payment.OrderID)),
			},
			"SK": {
				S: aws.String(fmt.Sprintf("PAYMENT#%s", payment.OrderID)),
			},
		},
		UpdateExpression: aws.String("SET #status = :status"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("Status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(payment.Status),
			},
		},
	}

	_, err := s.dynamoDB.UpdateItem(input)
	return err
}

func (s *paymentService) CreatePayment(ctx context.Context, event model.CreatedOrderEvent) error {
	paymentID := uuid.New().String()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.table),
		Item: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("PAYMENT#%s", event.OrderID)),
			},
			"SK": {
				S: aws.String(fmt.Sprintf("PAYMENT#%s", event.OrderID)),
			},
			"PaymentID": {
				S: aws.String(paymentID),
			},
			"OrderID": {
				S: aws.String(event.OrderID),
			},
			"Status": {
				S: aws.String("Incomplete"),
			},
		},
	}

	_, err := s.dynamoDB.PutItem(input)
	return err
}
