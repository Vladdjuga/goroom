package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"realTimeService/configuration"
)

// SimpleAuthMiddleware - simplified anonymous session middleware (no JWT required)
func SimpleAuthMiddleware(cfg *configuration.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate anonymous user ID from cookie or create new one
		sessionId := c.GetHeader("X-Session-ID")
		
		if sessionId == "" {
			// Create new anonymous session
			sessionId = uuid.New().String()
		}
		
		// Set the session ID for this request
		c.Set("user_sub", sessionId)
		c.Set("session_id", sessionId)
		
		// Proceed to the next handler
		c.Next()
	}
}
