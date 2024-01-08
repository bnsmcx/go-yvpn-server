package auth

import (
	"github.com/google/uuid"
	"log"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	deleteSession(r.Context().Value("id").(uuid.UUID))
	SetSessionCookie(w, "")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	pk := r.PostFormValue("portkey")

	account, err := Decrypt(pk)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid port-key", http.StatusBadRequest)
		return
	}

	createSession(account)
	SetSessionCookie(w, pk)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
