package db

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"yvpn_server/wg"
)

var database *gorm.DB

// Connect contains the startup and connection logic for the database
func Connect() error {
	dsn := "yvpn.db"
	d, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connecting: %s", err)
	}
	database = d

	// Migrate the schema
	err = database.AutoMigrate(&Account{}, &Endpoint{}, &Client{})
	if err != nil {
		return fmt.Errorf("schema automigration: %s", err)
	}
	return nil
}

func GetAccountByEmail(email string) (*Account, error) {
	var account Account
	result := database.Where("email = ?", email).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("record not found: %s", result.Error)
		}
		return nil, result.Error
	}
	return &account, nil
}

func GetAccount(id uuid.UUID) (*Account, error) {
	var account Account
	result := database.Preload("Endpoints").Where("id = ?", id).First(&account)
	if result.Error != nil {
		return nil, fmt.Errorf("record not found: %s", result.Error)
	}
	return &account, nil
}

func GetEndpoint(id int) (*Endpoint, error) {
	var endpoint Endpoint
	result := database.Where("id = ?", id).First(&endpoint)
	if result.Error != nil {
		return nil, fmt.Errorf("record not found: %s", result.Error)
	}
	return &endpoint, nil
}

func UpdateEndpointIPandClients(id int, ip string, clients map[string]wg.Keys) error {
	e, err := GetEndpoint(id)
	if err != nil {
		return err
	}

	e.IP = ip
	for k, v := range clients {
		err := e.AddClient(k, v.Private)
		if err != nil {
			return err
		}
	}
	return e.Save()
}
