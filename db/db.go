package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	err = database.AutoMigrate(&Account{}, &Endpoint{})
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
