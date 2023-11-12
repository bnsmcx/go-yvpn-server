package db

import (
	"fmt"
	"github.com/google/uuid"
)

// Account defines the Account record
type Account struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Activated         bool
	BearerToken       string
	DigitalOceanToken string
	Endpoints         []Endpoint `gorm:"foreignKey:AccountID"`
	Pin               string
}

func (a *Account) Activate() error {
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
