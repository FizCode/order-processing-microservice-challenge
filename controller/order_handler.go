package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func GetOrdersByCustomerIDHandler(orderService *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerID := r.URL.Query().Get("customer_id")
		if customerID == "" {
			http.Error(w, "Missing customer_id", http.StatusBadRequest)
			return
		}

		orders, err := orderService.GetOrdersByCustomerID(customerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}

func UpdateQtyByIDHandler(orderService *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		var updatedOrder model.Order
		if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = orderService.UpdateQtyByID(id, &updatedOrder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Order updated successfully")
	}
}

func DeleteOrderHandler(service *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		err = service.DeleteOrderByID(id)
		if err != nil {
			http.Error(w, "Failed to delete order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
