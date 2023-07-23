package main

import (
	"context"
	"encoding/json"
	"github.com/digitalocean/godo"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"strings"
	"yvpn_server/db"
	_ "yvpn_server/docs"
)

// @title yvpn API
// @version 0.0.1
// @description The public API for the yvpn server.
// @contact.name bnsmcx
// @contact.url https://github.com/bnsmcx/go-yvpn-server
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @license.name TODO
// @license.url TODO
// @host localhost:8000
// @BasePath /
func main() {
	r := chi.NewRouter()
	err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// serve docs
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/docs/doc.json"), //The url pointing to API definition
	))

	// API routes
	r.Get("/api/endpoint", CheckBearer(HandleGetEndpoints))
	r.Post("/api/endpoint", CheckBearer(HandleCreateEndpoint))

	http.ListenAndServe(":8000", r)
}

type Datacenter struct {
	Datacenter string `json:"datacenter"`
}

// HandleCreateEndpoint
// @Summary		Create a new endpoint
// @Description	create a new endpoint in the specified datacenter
// @Tags		endpoint
// @Security ApiKeyAuth
// @Accept 		json
// @Produce		text/plain
// @param token body Datacenter true "Datacenter"
// @Success		200 {string} string "new config"
// @Failure		400
// @Router		/api/endpoint [post]
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

// HandleGetEndpoints
// @Summary		Get all endpoints
// @Description	get all of a user's endpoints
// @Tags			endpoint
// @Security ApiKeyAuth
// @Produce		json
// @Success		200
// @Failure		400
// @Router			/api/endpoint [get]
func HandleGetEndpoints(w http.ResponseWriter, r *http.Request) {
	//TODO:  Get userID from Bearer token, return all their endpoints from db
}

// CheckBearer is middleware that validates the bearer token
func CheckBearer(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")

		// Split the header value by space
		headerParts := strings.Split(authHeader, " ")

		validFormat := len(headerParts) == 2 && strings.ToLower(headerParts[0]) == "bearer"
		if !validFormat {
			http.Error(w, "Invalid or missing Authorization header", http.StatusBadRequest)
			return
		}

		account, err := db.GetAccountByBearer(headerParts[1])
		if err != nil {
			http.Error(w, "Invalid or missing Authorization header", http.StatusBadRequest)
			return
		}

		// If all is good, call the next middleware or final handler
		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), "account", account.ID)))
	})
}
