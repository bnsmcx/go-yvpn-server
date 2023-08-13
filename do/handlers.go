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
		http.Error(w, "no a", http.StatusUnauthorized)
		return
	}

	a.DigitalOceanToken = r.PostFormValue("token")
	err = a.Save()

	http.Redirect(w, r, "/dash", http.StatusSeeOther)
}
