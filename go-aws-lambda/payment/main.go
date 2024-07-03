package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
	"payment/handler"
	"payment/service"
)

func main() {
	sess := session.Must(session.NewSession())
	paymentServ := service.NewPaymentService(
		dynamodb.New(sess),
		os.Getenv("PAYMENTS_TABLE"),
	)

	orderServ := service.NewOrderService(
		sqs.New(sess),
		os.Getenv("ORDERS_QUEUE_URL"),
	)

	handlerApp := handler.NewHandler(
		paymentServ,
		orderServ,
	)

	// Start the Lambda handler for ProcessPayment and CreatePayment
	lambda.Start(handlerApp.HandleRequest)
}
