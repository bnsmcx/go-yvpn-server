package ux

import (
	"github.com/google/uuid"
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
		InviteCode:        r.PostFormValue("invite-code"),
		DigitalOceanToken: r.PostFormValue("do-token"),
	}

	account, err := newAccount.Create()
	if err != nil {
		log.Println("Creating new account: " + err.Error())
		http.Error(w, "Invalid invite code", http.StatusBadRequest)
		return
	}

	log.Printf("Created new account: %s", account.ID.String())
	portkey, err := account.Encrypt()
	if err != nil {
		log.Println("Creating new account: " + err.Error())
		http.Error(w, "Invalid invite code", http.StatusBadRequest)
		return
	}
	RenderNewCreditNode(w, r, portkey)
}

func AddToken(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "parsing form", http.StatusBadRequest)
		return
	}

	a, err := auth.GetAccount(r.Context().Value("id").(uuid.UUID))
	if err != nil {
		log.Println(err)
		http.Error(w, "no account", http.StatusUnauthorized)
		return
	}

	a.DigitalOceanToken = r.PostFormValue("token")
	a.UpdateSessionStore()
	pk, err := a.Encrypt()
	if err != nil {
		log.Println(err)
	}
	auth.SetSessionCookie(w, pk)

	w.Header().Set("HX-Redirect", "/dashboard")
}
