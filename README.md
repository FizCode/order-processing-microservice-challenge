* run on terminal `docker-compose up --build`
* open postman -> POST on http://localhost:8080/orders with body:
```
{
"customer_id": "example-cust-id",
"product_id": "p-product_name",
"quantity": 1
}
```
* check on RabbitMQ http://localhost:15672/#/queues/%2F/orders -> click Get Message dropdown -> click button `Get Message(s)`

Project Structure:
```
order-processing-microservice-challenge/
├── controller/
│   └── order_controller.go
├── model/
│   └── model.go
├── mq/
│   └── publisher.go
├── repository/
│   └── order_repository.go
├── service/
│   └── order_service.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── init.sql
├── main.go
├── README.md
```