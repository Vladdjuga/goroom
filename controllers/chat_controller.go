package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ChatController handles the chat page
type ChatController struct{}

// NewChatController creates a new chat controller
func NewChatController() *ChatController {
	return &ChatController{}
}

// Index renders the chat page
func (c *ChatController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "chat.html", nil)
}
