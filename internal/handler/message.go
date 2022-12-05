package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"twitch_chat_analysis/internal/event"
	"twitch_chat_analysis/internal/model"
	"twitch_chat_analysis/internal/repository"
)

const (
	pathSaveMessage = "/message"
	pathReport      = "/message/list"
)

type handler struct {
	eventPublisher event.Publisher
	repo           repository.Repository
}

func NewHandler(eventPublisher event.Publisher, repo repository.Repository) *handler {
	return &handler{
		eventPublisher: eventPublisher,
		repo:           repo,
	}
}

func (h handler) RegisterRoutes(engine *gin.Engine) {
	engine.POST(pathSaveMessage, h.saveMessageHandler)
	engine.GET(pathReport, h.Report)
}

func (h handler) saveMessageHandler(c *gin.Context) {
	var msg model.Message

	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message body"})
		return
	}

	if !msg.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		return
	}

	if err := h.eventPublisher.Publish(context.Background(), &msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can not publish message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func (h handler) Report(c *gin.Context) {
	messages := h.repo.GetReports(c)
	c.JSON(http.StatusOK, gin.H{"data": messages})
}
