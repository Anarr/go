package event

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
	"twitch_chat_analysis/internal/model"
)

type Publisher interface {
	Publish(ctx context.Context, message *model.Message) error
}

type eventPublisher struct {
	Conn *amqp.Connection
}

func NewEventPublisher(conn *amqp.Connection) Publisher {
	return &eventPublisher{
		Conn: conn,
	}
}

func (ep eventPublisher) Publish(ctx context.Context, message *model.Message) error {
	ch, err := ep.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(topic, false, false, false, false, nil)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	if err != nil {
		return err
	}

	return nil
}
