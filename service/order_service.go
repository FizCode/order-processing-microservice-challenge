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
