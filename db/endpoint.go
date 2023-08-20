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
	ID         uuid.UUID `gorm:"primaryKey"`
	EndpointID int
	Config     string
	QR         []byte
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

	qr, err := wg.GenerateQR(config)
	if err != nil {
		return err
	}

	e.Clients = append(e.Clients, Client{
		ID:         uuid.New(),
		EndpointID: e.ID,
		Config:     config,
		QR:         qr,
	})

	return e.Save()
}
