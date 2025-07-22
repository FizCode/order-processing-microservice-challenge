package model

type Order struct {
	ID         int64  `json:"id"`
	CustomerID string `json:"customer_id"`
	ProductID  string `json:"product_id"`
	Quantity   int    `json:"quantity"`
}
