package db

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Account defines the Account record
type Account struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Activated         bool
	Email             string
	Password          []byte
	BearerToken       string
	DigitalOceanToken string
	Endpoints         []Endpoint `gorm:"foreignKey:AccountID"`
}

func (a *Account) Activate(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = hash
	a.Activated = true
	return a.Save()
}

func (a *Account) Save() error {
	result := database.Save(a)
	if result.Error != nil {
		return fmt.Errorf("saving account: %s", result.Error)
	}
	return nil
}
