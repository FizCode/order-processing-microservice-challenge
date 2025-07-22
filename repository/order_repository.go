package repository

import (
	"database/sql"

	"fizcode.dev/order-processing-microservice-challenge/model"
)

type OrderRepository struct {
	DB *sql.DB
}

func (r *OrderRepository) Save(order *model.Order) error {
	query := "INSERT INTO orders (customer_id, product_code, quantity) VALUES (?, ?, ?)"
	result, err := r.DB.Exec(query, order.CustomerID, order.ProductID, order.Quantity)
	if err != nil {
		return err
	}

	order.ID, _ = result.LastInsertId()
	return nil
}
