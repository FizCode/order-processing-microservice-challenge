package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	controller "fizcode.dev/order-processing-microservice-challenge/controller"
	"fizcode.dev/order-processing-microservice-challenge/mq"
	"fizcode.dev/order-processing-microservice-challenge/repository"
	"fizcode.dev/order-processing-microservice-challenge/service"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	rabbitUrl := os.Getenv("RABBITMQ_URL")
	queueName := "orders"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	ch.QueueDeclare(queueName, true, false, false, false, nil)

	orderRepository := &repository.OrderRepository{DB: db}
	publisher := &mq.Publisher{
		Channel: ch,
		Queue:   queueName,
	}
	orderService := &service.OrderService{
		OrderRepo: orderRepository,
		Publisher: publisher,
	}

	http.HandleFunc("/orders", controller.CreateOrderHandler(orderService))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
