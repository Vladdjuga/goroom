package models

import (
	"github.com/google/uuid"
	"time"
)

// Message represents a real-time chat message (no DB storage)
type Message struct {
	Type      string    `json:"type"`      // "message", "strangerJoined", "strangerLeft", etc
	Text      string    `json:"text"`
	UserId    uuid.UUID `json:"userId"`
	PairId    uuid.UUID `json:"pairId"`
	Timestamp time.Time `json:"timestamp"`
}

// NewMessage creates a new message
func NewMessage(text string, userId, pairId uuid.UUID) *Message {
	return &Message{
		Type:      "message",
		Text:      text,
		UserId:    userId,
		PairId:    pairId,
		Timestamp: time.Now(),
	}
}

// NewSystemMessage creates a system notification message
func NewSystemMessage(msgType string, pairId uuid.UUID) *Message {
	return &Message{
		Type:      msgType,
		PairId:    pairId,
		Timestamp: time.Now(),
	}
}
