package service

import (
	"context"
	"fmt"
	"order/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

// OrderService interface to handle queue operations
type OrderService interface {
	// CreateOrder creates a new order and sends a message to the queue
	CreateOrder(ctx context.Context, order model.CreateOrderRequest) (string, error)
	// CompleteOrder updates the order status
	CompleteOrder(ctx context.Context, orderID, status string) error
}

type orderService struct {
	dynamoDB *dynamodb.DynamoDB
	table    string
}

// NewOrderService creates a new OrderService
func NewOrderService(dynamoDB *dynamodb.DynamoDB, table string) OrderService {
	return &orderService{
		dynamoDB: dynamoDB,
		table:    table,
	}
}

func (os *orderService) CreateOrder(ctx context.Context, order model.CreateOrderRequest) (string, error) {
	orderID := uuid.New().String()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.table),
		Item: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("ORDER#%s", orderID)),
			},
			"SK": {
				S: aws.String(fmt.Sprintf("ORDER#%s", orderID)),
			},
			"UserID": {
				S: aws.String(order.UserID),
			},
			"Item": {
				S: aws.String(order.Item),
			},
			"Quantity": {
				N: aws.String(fmt.Sprintf("%d", order.Quantity)),
			},
			"TotalPrice": {
				N: aws.String(fmt.Sprintf("%d", order.TotalPrice)),
			},
			"Status": {
				S: aws.String("Pending"),
			},
		},
	}

	_, err := os.dynamoDB.PutItem(input)
	if err != nil {
		return "", err
	}

	return orderID, nil
}

func (os *orderService) CompleteOrder(ctx context.Context, orderID, status string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(os.table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("ORDER#%s", orderID)),
			},
			"SK": {
				S: aws.String(fmt.Sprintf("ORDER#%s", orderID)),
			},
		},
		UpdateExpression: aws.String("SET #status = :status"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("Status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(status),
			},
		},
	}

	_, err := os.dynamoDB.UpdateItem(input)
	return err
}
