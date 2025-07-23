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
	dsnNoDB := os.Getenv("MYSQL_DSN_NO_DB")
	dsnWithDB := os.Getenv("MYSQL_DSN")
	rabbitUrl := os.Getenv("RABBITMQ_URL")
	queueName := "orders"

	rootDB, err := sql.Open("mysql", dsnNoDB)
	if err != nil {
		log.Fatal(err)
	}
	defer rootDB.Close()

	_, err = rootDB.Exec("CREATE DATABASE IF NOT EXISTS orders")
	if err != nil {
		log.Fatal(err)
	}

	rootDB, err = sql.Open("mysql", dsnWithDB)
	if err != nil {
		log.Fatal(err)
	}
	defer rootDB.Close()

	if err := repository.InitDB(rootDB); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	ch.QueueDeclare(queueName, true, false, false, false, nil)

	orderRepository := &repository.OrderRepository{DB: rootDB}
	publisher := &mq.Publisher{
		Channel: ch,
		Queue:   queueName,
	}
	orderService := &service.OrderService{
		OrderRepo: orderRepository,
		Publisher: publisher,
	}

	http.HandleFunc("/orders", controller.CreateOrderHandler(orderService))
	http.HandleFunc("/orders/by-customer", controller.GetOrdersByCustomerIDHandler(orderService))
	http.HandleFunc("/orders/update", controller.UpdateQtyByIDHandler(orderService))
	http.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			controller.DeleteOrderHandler(orderService)(w, r)
			return
		}
		http.NotFound(w, r)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
