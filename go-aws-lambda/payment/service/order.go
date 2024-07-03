package service

import (
	"context"
	"encoding/json"
	"payment/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type OrderService interface {
	NotifyOrderCompletion(ctx context.Context, event model.OrderCompleteEvent) error
}

type orderService struct {
	sqs      *sqs.SQS
	queueURL string
}

func NewOrderService(sqs *sqs.SQS, queueURL string) OrderService {
	return &orderService{
		sqs:      sqs,
		queueURL: queueURL,
	}
}

func (s *orderService) NotifyOrderCompletion(ctx context.Context, event model.OrderCompleteEvent) error {
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
