package db

import (
	"fmt"
	"github.com/google/uuid"
)

// Account defines the Account record
type Account struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;"`
	Endpoints []Endpoint `gorm:"foreignKey:AccountID"`
}

func (a *Account) Save() error {
	result := database.Save(a)
	if result.Error != nil {
		return fmt.Errorf("saving account: %s", result.Error)
	}
	return nil
}
