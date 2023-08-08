package auth

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"yvpn_server/db"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
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

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	account, err := db.GetAccountByEmail(email)
	if err != nil {
		log.Println("failed login attempt, invalid email")
		http.Error(w, "invalid email", http.StatusUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword(account.Password, []byte(password))
	if err != nil {
		log.Println("failed login attempt, invalid password")
		http.Error(w, "invalid password", http.StatusUnauthorized)
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

	// Redirect to a protected page or return a success response
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
