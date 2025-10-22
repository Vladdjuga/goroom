package handlers

import (
	"realTimeService/interfaces"
	"realTimeService/models"

	"github.com/gin-gonic/gin"
)

// SendHandler implements MessageHandler interface
type SendHandler struct {
	container interfaces.Container
}

// NewSendHandler creates a new SendHandler instance
func NewSendHandler(container interfaces.Container) *SendHandler {
	return &SendHandler{
		container: container,
	}
}

// Handle processes the incoming message to send a message to stranger
func (h *SendHandler) Handle(ctx *gin.Context, client *models.Client,
	msg models.IncomingMessage, token string) error {
	
	// Validate the message
	if msg.Text == "" {
		ctx.JSON(400, gin.H{"error": "message text is required"})
		return nil
	}

	// Get the user's current pair
	pair, err := h.container.GetHub().MatchingService.GetPair(client.UserId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "you are not in an active chat"})
		return nil
	}

	if !pair.Active {
		ctx.JSON(400, gin.H{"error": "chat is not active"})
		return nil
	}

	// Create and send the message to partner
	outMsg := models.NewMessage(msg.Text, client.UserId, pair.ID)
	err = h.container.GetHub().SendMessageToPair(pair.ID, outMsg, client.UserId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to send message"})
		return err
	}

	return nil
}