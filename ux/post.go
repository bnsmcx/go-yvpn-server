package ux

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"yvpn_server/auth"
	"yvpn_server/db"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	newAccount := auth.NewCreditNode{
		InviteCode: r.PostFormValue("invite-code"),
	}

	account, err := newAccount.Create()
	if err != nil {
		log.Println("Creating new account: " + err.Error())
		http.Error(w, "Invalid invite code", http.StatusBadRequest)
		return
	}

	log.Printf("Created new account: %s", account.ID.String())
	RenderNewCreditNode(w, r, account.ID)
}

func ActivationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(r.PostFormValue("credit-code"))
	if err != nil {
		http.Error(w, "Invalid credit id", http.StatusBadRequest)
		return
	}

	if r.PostFormValue("password") != r.PostFormValue("confirm-password") {
		http.Error(w, "Passwords don't match", http.StatusBadRequest)
		return
	}

	a, err := db.GetAccount(id)
	if err != nil {
		http.Error(w, "Invalid credit id", http.StatusBadRequest)
		return
	}

	a.Activate(r.PostFormValue("password"))
}
