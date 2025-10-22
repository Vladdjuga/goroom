package handlers

import (
	"encoding/json"
	"realTimeService/interfaces"
	"realTimeService/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// NextStrangerHandler handles users skipping to the next stranger
type NextStrangerHandler struct {
	container interfaces.Container
}

// NewNextStrangerHandler creates a new NextStrangerHandler
func NewNextStrangerHandler(container interfaces.Container) *NextStrangerHandler {
	return &NextStrangerHandler{
		container: container,
	}
}

// Handle processes the next stranger request
func (h *NextStrangerHandler) Handle(ctx *gin.Context, client *models.Client,
	msg models.IncomingMessage, token string) error {

	hub := h.container.GetHub()

	// Get current pair
	currentPair, err := hub.MatchingService.GetPair(client.UserId)
	if err == nil && currentPair != nil {
		// Notify partner that user left
		partner := currentPair.GetPartner(client.UserId)
		if partner != nil {
			hub.NotifyStrangerLeft(partner.UserId)
		}

		// End current pair
		hub.MatchingService.EndPair(currentPair.ID)
	}

	// Try to find new match
	newPair, err := hub.MatchingService.FindMatch(client)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return err
	}

	// If no match yet, notify searching
	if newPair == nil {
		searchingMsg := models.NewSystemMessage(string(models.Searching), client.UserId)
		msgBytes, _ := json.Marshal(searchingMsg)
		client.Conn.WriteMessage(websocket.TextMessage, msgBytes)
		return nil
	}

	// New match found! Notify both users
	err = hub.NotifyStrangerJoined(newPair)
	if err != nil {
		return err
	}

	return nil
}
