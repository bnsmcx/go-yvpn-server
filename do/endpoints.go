package do

import (
	"context"
	"github.com/digitalocean/godo"
	"os"
	"yvpn_server/db"
)

type Datacenter struct {
	Datacenter string `json:"datacenter"`
}

func (d *Datacenter) CreateEndpoint() error {
	client := godo.NewFromToken(os.Getenv("DIGITAL_OCEAN_PAT"))
	ctx := context.TODO()

	createRequest := &godo.DropletCreateRequest{
		Name:   "yvpn-test",
		Region: d.Datacenter,
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
