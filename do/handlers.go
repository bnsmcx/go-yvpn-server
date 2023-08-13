package do

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"yvpn_server/db"
)

func AddToken(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "parsing form", http.StatusBadRequest)
		return
	}

	a, err := db.GetAccount(r.Context().Value("id").(uuid.UUID))
	if err != nil {
		log.Println(err)
		http.Error(w, "no account", http.StatusUnauthorized)
		return
	}

	a.DigitalOceanToken = r.PostFormValue("token")
	err = a.Save()
	if err != nil {
		log.Println(err)
		http.Error(w, "saving token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/dashboard")
}

func HandleCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "parsing form", http.StatusBadRequest)
		return
	}

	a, err := db.GetAccount(r.Context().Value("id").(uuid.UUID))
	if err != nil {
		log.Println(err)
		http.Error(w, "no account", http.StatusUnauthorized)
		return
	}

	endpoint := NewEndpoint{
		Token:      a.DigitalOceanToken,
		AccountID:  a.ID,
		Datacenter: r.FormValue("datacenter"),
	}

	err = endpoint.Create()
	if err != nil {
		log.Println(err)
		http.Error(w, "saving token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/dashboard")
}
