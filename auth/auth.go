package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"sync"
)

// sessionStore holds the user sessions
var sessionStore = make(map[uuid.UUID]Account)

// Account holds a user's details
type Account struct {
	DigitalOceanToken string
	ID                uuid.UUID
}

// Mutex to handle concurrent access to the sessionStore
var storeMutex sync.Mutex

// Create a session in the store
func createSession(id uuid.UUID) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	sessionID := hex.EncodeToString(bytes)
	storeMutex.Lock()
	sessionStore[sessionID] = id
	storeMutex.Unlock()
	return sessionID, nil
}

// Get a session from the store
func getSession(sessionID string) (uuid.UUID, bool) {
	storeMutex.Lock()
	session, found := sessionStore[sessionID]
	storeMutex.Unlock()
	return session, found
}

// Delete a session from the store
func deleteSession(sessionID string) {
	storeMutex.Lock()
	delete(sessionStore, sessionID)
	storeMutex.Unlock()
}

func GetAccount(id uuid.UUID) (*Account, error) {
	storeMutex.Lock()
	a, ok := sessionStore[id]
	storeMutex.Unlock()
	if !ok {
		return &Account{}, errors.New("Account not found")
	}
	return &a, nil
}
