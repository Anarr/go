package event

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"twitch_chat_analysis/internal/repository"
)

const topic = "report"

type Consumer interface {
	Consume() error
}

type eventConsumer struct {
	conn *amqp.Connection
	repo repository.Repository
}

func NewEventConsumer(conn *amqp.Connection, repo repository.Repository) Consumer {
	return &eventConsumer{
		conn: conn,
		repo: repo,
	}
}

func (ec eventConsumer) Consume() error {
	ch, err := ec.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(topic, false, false, false, false, nil)

	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			if err := ec.repo.Save(context.Background(), string(d.Body)); err != nil {
				log.Printf("Error happens while save event in db: %v\n", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages.")
	<-forever

	return nil
}
