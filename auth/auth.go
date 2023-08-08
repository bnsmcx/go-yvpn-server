package auth

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
	"sync"
)

// sessionStore holds the user sessions
var sessionStore = make(map[string]uuid.UUID)

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
