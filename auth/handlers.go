package auth

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"yvpn_server/db"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	deleteSession(r.Context().Value("session_id").(string))
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

	id := r.PostFormValue("credit-id")
	pin := r.PostFormValue("pin")

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid credit id", http.StatusBadRequest)
		return
	}
	account, err := db.GetAccount(uid)
	if err != nil {
		log.Println("failed login attempt, invalid email")
		http.Error(w, "invalid credit id", http.StatusUnauthorized)
	}

	if account.Pin != pin {
		log.Println("failed login attempt, invalid pin")
		http.Error(w, "invalid pin", http.StatusUnauthorized)
	}

	if !account.Activated {
		if err := account.Activate(); err != nil {
			log.Println("error activating card")
			http.Error(w, "", http.StatusInternalServerError)
		}
	}

	sessionID, err := createSession(account.ID)
	if err != nil {
		log.Println("error creating session key")
		http.Error(w, "", http.StatusInternalServerError)
	}

	// Create and set the cookie
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
