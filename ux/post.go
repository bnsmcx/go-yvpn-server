package ux

import (
	"log"
	"net/http"
	"yvpn_server/auth"
)

func PurchaseHandler(w http.ResponseWriter, r *http.Request) {
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
