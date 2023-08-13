package do

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/google/uuid"
	"yvpn_server/db"
)

type NewEndpoint struct {
	Token      string
	AccountID  uuid.UUID
	Datacenter string
}

func (e *NewEndpoint) Create() error {
	client := godo.NewFromToken(e.Token)
	ctx := context.TODO()

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

	dropletIP, err := droplet.PublicIPv4()
	if err != nil {
		return err
	}

	endpoint := db.Endpoint{
		ID:         droplet.ID,
		Datacenter: droplet.Region.Slug,
		IP:         dropletIP,
	}

	err = endpoint.Save()
	if err != nil {
		return err
	}
	return nil
}
