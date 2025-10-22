package hubs

import (
	"encoding/json"
	"fmt"
	"log"
	"realTimeService/models"
	"realTimeService/services"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MainHub struct {
	Clients         map[uuid.UUID]*models.Client
	MatchingService *services.MatchingService
	mut             sync.RWMutex
}

func NewMainHub() *MainHub {
	return &MainHub{
		Clients:         make(map[uuid.UUID]*models.Client),
		MatchingService: services.NewMatchingService(),
		mut:             sync.RWMutex{},
	}
}
func (h *MainHub) AddClient(client *models.Client) {
	h.mut.Lock()
	defer h.mut.Unlock()
	h.Clients[client.UserId] = client
	log.Printf("Client %s added to hub", client.UserId)
}

// SendMessageToPair sends a message to the partner in a pair
func (h *MainHub) SendMessageToPair(pairId uuid.UUID, message *models.Message, senderId uuid.UUID) error {
	pair, err := h.MatchingService.GetPairById(pairId)
	if err != nil {
		return fmt.Errorf("pair not found: %w", err)
	}

	if !pair.Active {
		return fmt.Errorf("pair is not active")
	}

	// Get the partner (not the sender)
	partner := pair.GetPartner(senderId)
	if partner == nil {
		return fmt.Errorf("partner not found")
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	err = partner.Conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		log.Printf("error sending message to client %s: %v", partner.UserId, err)
		return err
	}

	log.Printf("Message sent from %s to %s in pair %s", senderId, partner.UserId, pairId)
	return nil
}

// NotifyStrangerJoined notifies both users that they've been matched
func (h *MainHub) NotifyStrangerJoined(pair *models.ChatPair) error {
	notification := models.NewSystemMessage(string(models.StrangerJoined), pair.ID)
	
	msg1, _ := json.Marshal(notification)
	msg2, _ := json.Marshal(notification)

	err1 := pair.User1.Conn.WriteMessage(websocket.TextMessage, msg1)
	err2 := pair.User2.Conn.WriteMessage(websocket.TextMessage, msg2)

	if err1 != nil || err2 != nil {
		return fmt.Errorf("error notifying users: %v, %v", err1, err2)
	}

	log.Printf("Both users notified of match in pair %s", pair.ID)
	return nil
}

// NotifyStrangerLeft notifies a user that their partner has left
func (h *MainHub) NotifyStrangerLeft(userId uuid.UUID) error {
	h.mut.RLock()
	client, ok := h.Clients[userId]
	h.mut.RUnlock()

	if !ok {
		return fmt.Errorf("client not found")
	}

	notification := models.NewSystemMessage(string(models.StrangerLeft), uuid.Nil)
	messageBytes, _ := json.Marshal(notification)

	err := client.Conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		return fmt.Errorf("error notifying user: %w", err)
	}

	log.Printf("User %s notified that stranger left", userId)
	return nil
}

// RemoveClient removes a client from the hub and ends their pair if active
func (h *MainHub) RemoveClient(userId uuid.UUID) {
	h.mut.Lock()
	defer h.mut.Unlock()

	delete(h.Clients, userId)
	log.Printf("Client %s removed from hub", userId)

	// Try to get their pair and notify partner
	pair, err := h.MatchingService.GetPair(userId)
	if err == nil && pair.Active {
		// Notify partner
		partner := pair.GetPartner(userId)
		if partner != nil {
			go h.NotifyStrangerLeft(partner.UserId)
		}
		
		// End the pair
		h.MatchingService.EndPair(pair.ID)
	}

	// Also remove from waiting queue if they're there
	h.MatchingService.RemoveFromQueue(userId)
}

func (h *MainHub) GetClient(userId uuid.UUID) *models.Client {
	h.mut.RLock()
	defer h.mut.RUnlock()
	client, ok := h.Clients[userId]
	if !ok {
		return nil
	}
	log.Printf("Client %s retrieved from hub", userId)
	return client
}
