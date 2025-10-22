package models

import (
	"github.com/google/uuid"
	"time"
)

// ChatPair represents a matched pair of users chatting anonymously
type ChatPair struct {
	ID        uuid.UUID
	User1     *Client
	User2     *Client
	CreatedAt time.Time
	Active    bool
}

// NewChatPair creates a new chat pair between two clients
func NewChatPair(user1, user2 *Client) *ChatPair {
	return &ChatPair{
		ID:        uuid.New(),
		User1:     user1,
		User2:     user2,
		CreatedAt: time.Now(),
		Active:    true,
	}
}

// GetPartner returns the partner of the given user in the pair
func (cp *ChatPair) GetPartner(userId uuid.UUID) *Client {
	if cp.User1.UserId == userId {
		return cp.User2
	}
	return cp.User1
}

// HasUser checks if the given user is part of this pair
func (cp *ChatPair) HasUser(userId uuid.UUID) bool {
	return cp.User1.UserId == userId || cp.User2.UserId == userId
}

// Close marks the pair as inactive
func (cp *ChatPair) Close() {
	cp.Active = false
}
