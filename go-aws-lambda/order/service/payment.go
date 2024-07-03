package service

import (
	"context"
	"encoding/json"
	"order/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type PaymentService interface {
	NotifyPayment(ctx context.Context, event model.CreateOrderEvent) error
}

type paymentService struct {
	sqs      *sqs.SQS
	queueURL string
}

func NewPaymentService(sqs *sqs.SQS, queueURL string) PaymentService {
	return &paymentService{
		sqs:      sqs,
		queueURL: queueURL,
	}
}

func (s *paymentService) NotifyPayment(ctx context.Context, event model.CreateOrderEvent) error {
	eventBody, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = s.sqs.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(string(eventBody)),
	})

	return err
}
