package auth

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"sync"
)

// sessionStore holds the user sessions
var sessionStore = make(map[uuid.UUID]Account)

// Mutex to handle concurrent access to the sessionStore
var storeMutex sync.Mutex

// Account holds a user's details
type Account struct {
	DigitalOceanToken string
	ID                uuid.UUID
}

func Decrypt(pk string) (Account, error) {
	bytes, err := hex.DecodeString(pk)
	if err != nil {
		return Account{}, err
	}

	var a Account
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		return Account{}, err
	}
	return Account{}, nil
}

func (a *Account) Encrypt() (string, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Create a session in the store
func createSession(a Account) {
	storeMutex.Lock()
	sessionStore[a.ID] = a
	storeMutex.Unlock()
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
