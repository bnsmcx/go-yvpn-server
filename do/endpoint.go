package do

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"time"
	"yvpn_server/db"
	"yvpn_server/wg"
)

type NewEndpoint struct {
	Token      string
	AccountID  uuid.UUID
	Datacenter string
}

func (e *NewEndpoint) Create() error {
	client := godo.NewFromToken(e.Token)
	ctx := context.TODO()

	type KeyPair struct {
		PubKey  wgtypes.Key
		PrivKey wgtypes.Key
	}

	pub, priv, err := wg.GenerateKeys()
	if err != nil {
		return err
	}
	serverKeys := KeyPair{
		PubKey:  pub,
		PrivKey: priv,
	}

	var clientKeys = make(map[string]KeyPair)

	for i := 2; i <= 255; i++ {
		pub, priv, err := wg.GenerateKeys()
		if err != nil {
			return err
		}
		keys := KeyPair{
			PubKey:  pub,
			PrivKey: priv,
		}
		clientKeys[fmt.Sprintf("10.0.0.%d", i)] = keys
	}

	fmt.Println(serverKeys)
	for k, v := range clientKeys {
		fmt.Println(k)
		fmt.Println(v)
	}
	//serverConfig, err := wg.GenerateServerConfig(serverKeys, clientKeys)

	createRequest := &godo.DropletCreateRequest{
		Name:   "yvpn-test",
		Region: e.Datacenter,
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			ID: 110391971,
		},
	}
	droplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		return err
	}

	go awaitIP(e.Token, droplet.ID)

	endpoint := db.Endpoint{
		ID:         droplet.ID,
		Datacenter: droplet.Region.Slug,
		AccountID:  e.AccountID,
	}

	err = endpoint.Save()
	if err != nil {
		return err
	}
	return nil
}

func awaitIP(token string, id int) {
	for i := 0; i < 36; i++ {
		time.Sleep(10 * time.Second)
		client := godo.NewFromToken(token)
		droplet, _, err := client.Droplets.Get(context.TODO(), id)
		if err != nil {
			log.Println(err)
			return
		} else if droplet.Status != "active" {
			continue
		}

		ip, err := droplet.PublicIPv4()
		if err != nil {
			log.Println(err)
			return
		}
		err = db.UpdateEndpointIP(id, ip)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
