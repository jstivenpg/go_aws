package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"order/handler"
	"order/service"

	"os"
)

func main() {
	sess := session.Must(session.NewSession())
	orderServ := service.NewOrderService(
		dynamodb.New(sess),
		os.Getenv("ORDERS_TABLE"),
	)

	paymentServ := service.NewPaymentService(
		sqs.New(sess),
		os.Getenv("PAYMENTS_QUEUE_URL"),
	)

	handlerApp := handler.NewHandler(
		orderServ,
		paymentServ,
	)

	lambda.Start(handlerApp.HandleRequest)
}
