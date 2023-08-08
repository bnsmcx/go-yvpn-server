package ux

import (
	"log"
	"net/http"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		log.Println("What?")

		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		confirmPassword := r.PostFormValue("confirm-password")
		inviteCode := r.PostFormValue("invite-code")

		log.Println("Email:", email)
		log.Println("Password:", password)
		log.Println("Confirm Password:", confirmPassword)
		log.Println("Invite Code:", inviteCode)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
