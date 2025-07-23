package mq

import (
	"context"
	"encoding/json"

	"fizcode.dev/order-processing-microservice-challenge/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Channel *amqp.Channel
	Queue   string
}

func (p *Publisher) PublishOrder(order *model.Order) error {
	body, err := json.Marshal(order)
	ctx := context.Background()
	if err != nil {
		return err
	}
	return p.Channel.PublishWithContext(
		ctx,
		"",
		p.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) PublishOrderDeleted(orderID int) error {
	payload := map[string]interface{}{
		"type": "deleted",
		"id":   orderID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return p.Channel.PublishWithContext(
		ctx,
		"",
		p.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
