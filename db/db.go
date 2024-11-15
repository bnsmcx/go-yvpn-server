package db

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"yvpn_server/wg"
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
	err = database.AutoMigrate(&Account{}, &Endpoint{}, &Client{})
	if err != nil {
		return fmt.Errorf("schema automigration: %s", err)
	}
	return nil
}

func GetAccount(id uuid.UUID) (*Account, error) {
	var account Account
	result := database.Preload("Endpoints.Clients").Where("id = ?", id).First(&account)
	if result.Error != nil {
		return nil, fmt.Errorf("record not found: %s", result.Error)
	}
	return &account, nil
}

func GetEndpoint(id int) (*Endpoint, error) {
	var endpoint Endpoint
	result := database.Preload("Clients").Where("id = ?", id).First(&endpoint)
	if result.Error != nil {
		return nil, fmt.Errorf("record not found: %s", result.Error)
	}
	return &endpoint, nil
}

func GetClient(id uuid.UUID) (*Client, error) {
	var client Client
	result := database.Where("id = ?", id).First(&client)
	if result.Error != nil {
		return nil, fmt.Errorf("record not found: %s", result.Error)
	}
	return &client, nil
}

func UpdateEndpointIPandClients(id int, ip string, clients map[string]wg.Keys) error {
	for k, v := range clients {
		e, err := GetEndpoint(id) // make sure we get an updated obj after each loop's write
		if err != nil {
			return err
		}

		e.IP = ip
		err = e.AddClient(k, v.Private)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetClientConfigFile(w http.ResponseWriter, r *http.Request) {
	endpointID, err := strconv.Atoi(chi.URLParam(r, "endpoint"))
	if err != nil {
		return
	}

	clientID := uuid.MustParse(chi.URLParam(r, "client"))
	e, err := GetEndpoint(endpointID)
	if err != nil {
		return
	}

	c, err := e.GetClient(clientID)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s.conf", e.Datacenter))
	_, err = w.Write([]byte(c.Config))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func HandleRenameClient(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	id := uuid.MustParse(chi.URLParam(r, "id"))

	c, err := GetClient(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "error renaming client", http.StatusBadRequest)
		return
	}

	c.Name = name
	if err := c.Save(); err != nil {
		log.Println(err)
		http.Error(w, "error renaming client", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func HandleDeleteClient(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	c, err := GetClient(id)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		log.Println(err)
		return
	}
	if err := c.Delete(); err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		log.Println(err)
		return
	}
	w.Header().Set("HX-Redirect", "/dashboard")
}
