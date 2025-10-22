package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HomeController handles the home page
type HomeController struct{}

// NewHomeController creates a new home controller
func NewHomeController() *HomeController {
	return &HomeController{}
}

// Index renders the home page
func (c *HomeController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home.html", nil)
}
