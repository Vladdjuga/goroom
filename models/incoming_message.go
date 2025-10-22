package models

import "github.com/google/uuid"

type MessageType string

const (
	// User actions
	FindMatch       MessageType = "findMatch"       // Find a random stranger
	SendMessage     MessageType = "sendMessage"     // Send message to stranger
	NextStranger    MessageType = "nextStranger"    // Skip to next stranger
	StopChat        MessageType = "stopChat"        // Stop chatting
	Typing          MessageType = "typing"          // User is typing notification
	
	// System notifications (outgoing)
	StrangerJoined  MessageType = "strangerJoined"  // Stranger connected
	StrangerLeft    MessageType = "strangerLeft"    // Stranger disconnected
	Searching       MessageType = "searching"       // Looking for stranger
)

type IncomingMessage struct {
	Type   MessageType `json:"type"`
	PairId uuid.UUID   `json:"pairId,omitempty"` // Optional: current pair ID
	Text   string      `json:"text,omitempty"`   // Optional: message text
}
