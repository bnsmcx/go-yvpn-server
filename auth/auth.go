package auth

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
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
	return a, nil
}

func (a *Account) Encrypt() (string, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a *Account) UpdateSessionStore() {
	storeMutex.Lock()
	sessionStore[a.ID] = *a
	storeMutex.Unlock()
}

// Create a session in the store
func createSession(a Account) {
	storeMutex.Lock()
	sessionStore[a.ID] = a
	storeMutex.Unlock()
}

// Get a session from the store
func getSession(pk string) (*Account, error) {
	a, err := Decrypt(pk)
	if err != nil {
		return &Account{}, err
	}
	storeMutex.Lock()
	a, found := sessionStore[a.ID]
	storeMutex.Unlock()
	if !found {
		return &Account{}, errors.New("session not found")
	}
	return &a, nil
}

// Delete a session from the store
func deleteSession(id uuid.UUID) {
	storeMutex.Lock()
	delete(sessionStore, id)
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

func SetSessionCookie(w http.ResponseWriter, pk string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    pk,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}
