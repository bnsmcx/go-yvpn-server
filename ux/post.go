package ux

import (
	"log"
	"net/http"
	"yvpn_server/auth"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	newUser := auth.NewUser{
		Email:       r.PostFormValue("email"),
		Password:    r.PostFormValue("password"),
		ConfirmPass: r.PostFormValue("confirm-password"),
		InviteCode:  r.PostFormValue("invite-code"),
	}

	account, err := auth.CreateAccount(newUser)
	if err != nil {
		log.Println("Creating new account: " + err.Error())
	}

	log.Println("Created new account: <TODO>")
}
