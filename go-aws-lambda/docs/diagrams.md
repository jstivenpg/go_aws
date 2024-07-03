Create Order

```bash
@startuml
actor Client
participant "OrderService" as OS
participant "DynamoDB" as DB
participant "PaymentQueue" as PQ

Client -> OS : POST /order
OS -> OS : CreateOrder(request)
OS -> DB : PutItem(CreateOrderRequest)
DB --> OS : Confirmation
OS -> PQ : SendMessage(CreateOrderEvent)
PQ --> OS : Confirmation
OS --> Client : HTTP 200 OK
@enduml
```

Complete Order:
```bash
@startuml
participant "OrderQueue" as OQ
participant "OrderService" as OS
participant "DynamoDB" as DB

OQ -> OS : SQS Event with orderId and status
OS -> OS : CompleteOrder(orderId, status)
OS -> DB : UpdateItem(orderId, status)
DB --> OS : Confirmation
@enduml
```

### Order Payment

Create Payment

```bash
@startuml
participant "PaymentQueue" as PQ
participant "PaymentService" as PS
participant "DynamoDB" as DB

PQ -> PS : SQS Event with order details
PS -> PS : CreatePayment(event)
PS -> DB : PutItem(CreatePaymentRequest)
DB --> PS : Confirmation
@enduml
```

Process payment

```bash
@startuml
actor Client
participant "PaymentService" as PS
participant "DynamoDB" as DB
participant "OrderQueue" as OQ

Client -> PS : POST /payment
PS -> PS : ProcessPayment(request)
PS -> DB : UpdateItem(ProcessPaymentRequest)
DB --> PS : Confirmation
PS -> OQ : SendMessage(OrderCompletedEvent)
OQ --> PS : Confirmation
PS --> Client : HTTP 200 OK
@enduml
```