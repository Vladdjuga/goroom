package handlers

import (
	"realTimeService/interfaces"
	"realTimeService/models"

	"github.com/gin-gonic/gin"
)

// StopChatHandler handles users stopping their chat
type StopChatHandler struct {
	container interfaces.Container
}

// NewStopChatHandler creates a new StopChatHandler
func NewStopChatHandler(container interfaces.Container) *StopChatHandler {
	return &StopChatHandler{
		container: container,
	}
}

// Handle processes the stop chat request
func (h *StopChatHandler) Handle(ctx *gin.Context, client *models.Client,
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

	// Remove from waiting queue if there
	hub.MatchingService.RemoveFromQueue(client.UserId)

	return nil
}
