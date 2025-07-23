package controller

import (
	"encoding/json"
	"net/http"

	"fizcode.dev/order-processing-microservice-challenge/model"
	"fizcode.dev/order-processing-microservice-challenge/service"
)

func CreateOrderHandler(orderService *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order model.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if order.CustomerID == "" || order.ProductID == "" || order.Quantity <= 0 {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		if err := orderService.CreateOrder(&order); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}
