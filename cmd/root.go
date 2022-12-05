package cmd

import (
	"github.com/gin-gonic/gin"
	"log"
	"twitch_chat_analysis/config"
	"twitch_chat_analysis/internal/event"
	"twitch_chat_analysis/internal/handler"
	"twitch_chat_analysis/internal/repository"
)

func Execute() {
	appConfig, err := config.InitAppConfig()
	failOnErr(err)

	amqp, err := newAMQP(appConfig)
	failOnErr(err)
	defer amqp.Close()

	redisClient := newRedisClient(appConfig)
	repo := repository.NewRepository(redisClient)

	eventConsumer := event.NewEventConsumer(amqp, repo)

	go func() {
		err := eventConsumer.Consume()
		failOnErr(err)
	}()

	eventPublisher := event.NewEventPublisher(amqp)

	h := handler.NewHandler(eventPublisher, repo)

	r := gin.Default()
	h.RegisterRoutes(r)
	err = r.Run()
	failOnErr(err)
}

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
