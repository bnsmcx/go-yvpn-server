package db

import (
	"fmt"
	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"yvpn_server/wg"
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
	Config     string
}

func (e *Endpoint) Save() error {
	result := database.Save(e)
	if result.Error != nil {
		return fmt.Errorf("creating endpoint db record: %s", result.Error)
	}
	return nil
}

func (e *Endpoint) AddClient(clientIP string, privKey wgtypes.Key) error {
	config, err := wg.GenerateClientConfig(e.IP, clientIP, e.PublicKey, privKey.String())
	if err != nil {
		return err
	}

	e.Clients = append(e.Clients, Client{
		ID:         "",
		EndpointID: e.ID,
		Config:     config,
	})

	return e.Save()
}
