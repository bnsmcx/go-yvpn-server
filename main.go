package main

import (
	"context"
	"encoding/json"
	"github.com/digitalocean/godo"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"yvpn_server/auth"
	"yvpn_server/db"
	"yvpn_server/ux"
)

func main() {
	r := chi.NewRouter()
	err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Get("/", ux.RenderLanding)
		r.Get("/signup", ux.RenderSignup)
		r.Get("/login", ux.RenderLogin)
		r.Get("/logout", auth.HandleLogout)
	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(auth.UserSession)
		r.Get("/manage", ux.RenderManage)
	})

	log.Fatalln(http.ListenAndServe(":8000", r))
}

type Datacenter struct {
	Datacenter string `json:"datacenter"`
}

func HandleCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	var dc Datacenter
	err := json.NewDecoder(r.Body).Decode(&dc)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid datacenter json", http.StatusBadRequest)
		return
	}

	client := godo.NewFromToken(os.Getenv("DIGITAL_OCEAN_PAT"))
	ctx := context.TODO()

	createRequest := &godo.DropletCreateRequest{
		Name:   "yvpn-test",
		Region: dc.Datacenter,
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			ID: 110391971,
		},
	}
	d, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, "error creating droplet", http.StatusInternalServerError)
		return
	}

	dropletIP, err := d.PublicIPv4()
	if err != nil {
		log.Println(err)
		http.Error(w, "getting droplet ip", http.StatusInternalServerError)
		return
	}

	endpoint := db.Endpoint{
		ID:         d.ID,
		Datacenter: d.Region.Slug,
		AccountID:  r.Context().Value("account").(uuid.UUID),
		IP:         dropletIP,
	}

	err = endpoint.Save()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
