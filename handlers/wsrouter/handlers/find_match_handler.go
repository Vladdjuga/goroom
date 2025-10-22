package handlers

import (
	"encoding/json"
	"realTimeService/interfaces"
	"realTimeService/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// FindMatchHandler handles users looking for a random stranger to chat with
type FindMatchHandler struct {
	container interfaces.Container
}

// NewFindMatchHandler creates a new FindMatchHandler
func NewFindMatchHandler(container interfaces.Container) *FindMatchHandler {
	return &FindMatchHandler{
		container: container,
	}
}

// Handle processes the find match request
func (h *FindMatchHandler) Handle(ctx *gin.Context, client *models.Client,
	msg models.IncomingMessage, token string) error {

	hub := h.container.GetHub()
	
	// Try to find a match
	pair, err := hub.MatchingService.FindMatch(client)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return err
	}

	// If no match yet (added to queue), notify user they're searching
	if pair == nil {
		searchingMsg := models.NewSystemMessage(string(models.Searching), client.UserId)
		msgBytes, _ := json.Marshal(searchingMsg)
		client.Conn.WriteMessage(websocket.TextMessage, msgBytes)
		return nil
	}

	// Match found! Notify both users
	err = hub.NotifyStrangerJoined(pair)
	if err != nil {
		return err
	}

	return nil
}
