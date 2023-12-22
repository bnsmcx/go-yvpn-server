package db

import (
	"errors"
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
	Active     bool
	EndpointID int
	Config     string
	QR         string
}

func (c *Client) Save() error {
	result := database.Save(c)
	if result.Error != nil {
		return fmt.Errorf("saving client config: %s", result.Error)
	}
	return nil
}

func (c *Client) Delete() error {
	result := database.Delete(c)
	return result.Error
}

func (e *Endpoint) Save() error {
	result := database.Save(e)
	if result.Error != nil {
		return fmt.Errorf("creating endpoint db record: %s", result.Error)
	}
	return nil
}

func (e *Endpoint) GetClient(id uuid.UUID) (*Client, error) {
	for _, c := range e.Clients {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, errors.New("client not found")
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

func (e *Endpoint) Delete() error {
	result := database.Delete(e)
	return result.Error
}

func (e *Endpoint) DeleteClientConfigsForEndpoint() {
	for _, c := range e.Clients {
		database.Delete(c)
	}
}

func (e *Endpoint) ActivateClient() error {
	for _, c := range e.Clients {
		if !c.Active {
			c.Active = true
			return c.Save()
		}
	}
	return nil
}
