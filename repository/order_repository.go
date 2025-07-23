package repository

import (
	"database/sql"
	"errors"
	"fizcode.dev/order-processing-microservice-challenge/model"
	"fmt"
)

type OrderRepository struct {
	DB *sql.DB
}

func (r *OrderRepository) Save(order *model.Order) error {
	insertQuery := "INSERT INTO orders (customer_id, product_id, quantity) VALUES (?, ?, ?)"
	selectQuery := "SELECT id, customer_id, product_id, quantity, created_at, updated_at FROM orders WHERE id = ?"
	result, err := r.DB.Exec(insertQuery, order.CustomerID, order.ProductID, order.Quantity)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = id

	row := r.DB.QueryRow(selectQuery, id)
	err = row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.ProductID,
		&order.Quantity,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("order with id %d not found after insert", id)
	} else if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) FindByID(id int) (*model.Order, error) {
	row := r.DB.QueryRow("SELECT id, customer_id, product_id, quantity, created_at, updated_at FROM orders WHERE id = ?", id)

	var order model.Order
	err := row.Scan(&order.ID, &order.CustomerID, &order.ProductID, &order.Quantity, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) FindByCustomerID(customerID string) ([]model.Order, error) {
	query := "SELECT id, customer_id, product_id, quantity, created_at, updated_at FROM orders WHERE customer_id = ?"
	rows, err := r.DB.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.ProductID,
			&order.Quantity,
			&order.CreatedAt,
			&order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateQtyByID(ID int, updatedOrder *model.Order) error {
	query := "UPDATE orders SET quantity = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.DB.Exec(query, updatedOrder.Quantity, ID)

	return err
}

func (r *OrderRepository) DeleteByID(id int) error {
	_, err := r.DB.Exec("DELETE FROM orders WHERE id = ?", id)
	return err
}
