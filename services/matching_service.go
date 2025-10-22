package services

import (
	"fmt"
	"log"
	"realTimeService/models"
	"sync"

	"github.com/google/uuid"
)

// MatchingService handles pairing users for anonymous chats
type MatchingService struct {
	waitingQueue []*models.Client
	activePairs  map[uuid.UUID]*models.ChatPair
	userToPair   map[uuid.UUID]uuid.UUID // userId -> pairId mapping
	mu           sync.RWMutex
}

// NewMatchingService creates a new matching service
func NewMatchingService() *MatchingService {
	return &MatchingService{
		waitingQueue: make([]*models.Client, 0),
		activePairs:  make(map[uuid.UUID]*models.ChatPair),
		userToPair:   make(map[uuid.UUID]uuid.UUID),
		mu:           sync.RWMutex{},
	}
}

// FindMatch tries to find a partner for the given client
// Returns the created pair if match found, nil if added to queue
func (m *MatchingService) FindMatch(client *models.Client) (*models.ChatPair, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if user is already in a pair
	if pairId, exists := m.userToPair[client.UserId]; exists {
		if pair, ok := m.activePairs[pairId]; ok && pair.Active {
			return nil, fmt.Errorf("user already in active chat")
		}
	}

	// If no one is waiting, add to queue
	if len(m.waitingQueue) == 0 {
		m.waitingQueue = append(m.waitingQueue, client)
		log.Printf("User %s added to waiting queue", client.UserId)
		return nil, nil // nil means waiting for match
	}

	// Take first person from queue
	stranger := m.waitingQueue[0]
	m.waitingQueue = m.waitingQueue[1:]

	// Don't match with yourself
	if stranger.UserId == client.UserId {
		m.waitingQueue = append(m.waitingQueue, client)
		return nil, nil
	}

	// Create pair
	pair := models.NewChatPair(client, stranger)
	m.activePairs[pair.ID] = pair
	m.userToPair[client.UserId] = pair.ID
	m.userToPair[stranger.UserId] = pair.ID

	log.Printf("Matched users %s and %s in pair %s", client.UserId, stranger.UserId, pair.ID)
	return pair, nil
}

// RemoveFromQueue removes a client from the waiting queue
func (m *MatchingService) RemoveFromQueue(userId uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, client := range m.waitingQueue {
		if client.UserId == userId {
			m.waitingQueue = append(m.waitingQueue[:i], m.waitingQueue[i+1:]...)
			log.Printf("User %s removed from waiting queue", userId)
			return
		}
	}
}

// GetPair returns the active pair for the given user
func (m *MatchingService) GetPair(userId uuid.UUID) (*models.ChatPair, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pairId, exists := m.userToPair[userId]
	if !exists {
		return nil, fmt.Errorf("user not in any pair")
	}

	pair, ok := m.activePairs[pairId]
	if !ok {
		return nil, fmt.Errorf("pair not found")
	}

	return pair, nil
}

// GetPairById returns a pair by its ID
func (m *MatchingService) GetPairById(pairId uuid.UUID) (*models.ChatPair, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pair, ok := m.activePairs[pairId]
	if !ok {
		return nil, fmt.Errorf("pair not found")
	}

	return pair, nil
}

// EndPair closes a pair and removes both users from mapping
func (m *MatchingService) EndPair(pairId uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pair, ok := m.activePairs[pairId]
	if !ok {
		return fmt.Errorf("pair not found")
	}

	pair.Close()
	delete(m.userToPair, pair.User1.UserId)
	delete(m.userToPair, pair.User2.UserId)
	delete(m.activePairs, pairId)

	log.Printf("Pair %s ended", pairId)
	return nil
}

// EndUserPair ends the pair that the user is currently in
func (m *MatchingService) EndUserPair(userId uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pairId, exists := m.userToPair[userId]
	if !exists {
		return fmt.Errorf("user not in any pair")
	}

	pair, ok := m.activePairs[pairId]
	if !ok {
		return fmt.Errorf("pair not found")
	}

	pair.Close()
	delete(m.userToPair, pair.User1.UserId)
	delete(m.userToPair, pair.User2.UserId)
	delete(m.activePairs, pairId)

	log.Printf("Pair %s ended by user %s", pairId, userId)
	return nil
}

// GetQueueSize returns the number of users waiting for a match
func (m *MatchingService) GetQueueSize() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.waitingQueue)
}

// GetActivePairsCount returns the number of active pairs
func (m *MatchingService) GetActivePairsCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.activePairs)
}
