package db

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB

// Account defines the Account record
type Account struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;"`
	BearerToken       string
	DigitalOceanToken string
	Endpoints         []Endpoint `gorm:"foreignKey:AccountID"`
}

// Endpoint defines the Endpoint record
type Endpoint struct {
	ID         int `gorm:"primaryKey"`
	Datacenter string
	AccountID  uuid.UUID
	IP         string
}

func (e *Endpoint) Save() error {
	result := database.Create(e)
	if result.Error != nil {
		return fmt.Errorf("creating endpoint db record: %s", result.Error)
	}
	return nil
}

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

func GetAccountByBearer(bearer string) (*Account, error) {
	var account = Account{BearerToken: bearer}
	result := database.First(&account)
	if result.Error != nil {
		return nil, fmt.Errorf("record not found: %s", result.Error)
	}
	return &account, nil
}
