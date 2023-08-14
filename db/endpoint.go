package db

import (
	"fmt"
	"github.com/google/uuid"
)

// Endpoint defines the Endpoint record
type Endpoint struct {
	ID         int `gorm:"primaryKey"`
	Datacenter string
	AccountID  uuid.UUID
	IP         string
	PublicKey  string
	PrivateKey string
	Clients    []Client `gorm:"foreignKey:EndpointID"`
}

type Client struct {
	ID         string `gorm:"primaryKey"`
	EndpointID int
	IP         string
	PublicKey  string
	PrivateKey string
}

func (e *Endpoint) Save() error {
	result := database.Save(e)
	if result.Error != nil {
		return fmt.Errorf("creating endpoint db record: %s", result.Error)
	}
	return nil
}
