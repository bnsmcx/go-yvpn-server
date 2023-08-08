package db

import "github.com/google/uuid"

// Account defines the Account record
type Account struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Email             string
	Password          string
	BearerToken       string
	DigitalOceanToken string
	Endpoints         []Endpoint `gorm:"foreignKey:AccountID"`
}
