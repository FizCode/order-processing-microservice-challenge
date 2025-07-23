package service

import (
	"fizcode.dev/order-processing-microservice-challenge/model"
	"fizcode.dev/order-processing-microservice-challenge/mq"
	"fizcode.dev/order-processing-microservice-challenge/repository"
)

type OrderService struct {
	OrderRepo *repository.OrderRepository
	Publisher *mq.Publisher
}

func (s *OrderService) CreateOrder(order *model.Order) error {
	if err := s.OrderRepo.Save(order); err != nil {
		return err
	}
	return s.Publisher.PublishOrder(order)
}

func (s *OrderService) GetOrdersByCustomerID(customerID string) ([]model.Order, error) {
	return s.OrderRepo.FindByCustomerID(customerID)
}

func (s *OrderService) UpdateQtyByID(ID int, updatedOrder *model.Order) error {
	if err := s.OrderRepo.UpdateQtyByID(ID, updatedOrder); err != nil {
		return err
	}

	updated, err := s.OrderRepo.FindByID(ID)
	if err != nil {
		return err
	}
	return s.Publisher.PublishOrder(updated)
}

func (s *OrderService) DeleteOrderByID(id int) error {
	err := s.OrderRepo.DeleteByID(id)
	if err != nil {
		return err
	}
	return s.Publisher.PublishOrderDeleted(id)
}
