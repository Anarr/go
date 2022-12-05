package cmd

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"twitch_chat_analysis/config"
)

func newAMQP(conf *config.AppConfig) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.RabbitMQ.Username, conf.RabbitMQ.Password, conf.RabbitMQ.Host, conf.RabbitMQ.Port))
}

func newRedisClient(conf *config.AppConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%d", conf.Redis.Port),
		Password: "",
		DB:       0,
	})
}
