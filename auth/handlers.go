package auth

import (
	"github.com/google/uuid"
	"log"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	deleteSession(uuid.MustParse(r.Context().Value("id").(string)))
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

	// Create and set the cookie
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    pk,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
