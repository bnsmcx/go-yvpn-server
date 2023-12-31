package ux

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"yvpn_server/auth"
	"yvpn_server/db"
	"yvpn_server/do"
)

type pageData struct {
	LoggedIn  bool
	Account   *db.Account
	PageTitle string
	PortKey   string
	UserData  *auth.Account
}

func RenderLanding(w http.ResponseWriter, r *http.Request) {
	// Parse the templates
	tmpl, err := template.ParseFiles("templates/base.html", "templates/landing.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if user is logged in, add to page data
	var pd pageData
	pd.PageTitle = "yVPN"
	id := r.Context().Value("id")
	if id != nil {
		pd.LoggedIn = true
	}

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderPurchase(w http.ResponseWriter, r *http.Request) {
	// Parse the templates
	tmpl, err := template.ParseFiles("templates/base.html", "templates/purchase.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if user is logged in, add to page data
	var pd pageData
	pd.PageTitle = "Purchase"
	id := r.Context().Value("id")
	if id != nil {
		pd.LoggedIn = true
	}

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderLogin(w http.ResponseWriter, r *http.Request) {
	// Parse the templates
	tmpl, err := template.ParseFiles("templates/base.html", "templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if user is logged in, add to page data
	var pd pageData
	pd.PageTitle = "Login"
	id := r.Context().Value("id")
	if id != nil {
		pd.LoggedIn = true
	}

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderDashboard(w http.ResponseWriter, r *http.Request) {
	// Get the user details
	ud, err := auth.GetAccount(r.Context().Value("id").(uuid.UUID))

	// Get the user db record
	dbRecord, err := db.GetAccount(ud.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// populate page data
	pd := pageData{
		LoggedIn:  true,
		PageTitle: "Dashboard",
		UserData:  ud,
		Account:   dbRecord,
	}

	// Parse the templates
	tmpl, err := template.ParseFiles("templates/base.html", "templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderAddToken(w http.ResponseWriter, r *http.Request) {
	// Parse the templates
	tmpl, err := template.ParseFiles("templates/add-token.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderAddEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get account
	a, err := auth.GetAccount(r.Context().Value("id").(uuid.UUID))
	if err != nil {
		log.Println(err)
		http.Error(w, "getting account", http.StatusUnauthorized)
		return
	}

	// Get available datacenters
	dc, err := do.GetDatacenters(a.DigitalOceanToken)
	if err != nil {
		log.Println(err)
		http.Error(w, "getting datacenters", http.StatusInternalServerError)
		return
	}

	// Parse the templates
	tmpl, err := template.ParseFiles("templates/add-endpoint.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the "layout" template and send it to the ResponseWriter.
	data := struct {
		Datacenters []string
	}{
		Datacenters: dc,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderNewCreditNode(w http.ResponseWriter, r *http.Request, portkey string) {
	// Parse the templates
	tmpl, err := template.ParseFiles("templates/base.html", "templates/credit_node.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// populate page data
	pd := pageData{
		LoggedIn:  r.Context().Value("id") != nil,
		PageTitle: "New Credit",
		PortKey:   portkey,
	}

	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ActivateClient(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "endpoint"))
	e, err := db.GetEndpoint(id)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		log.Println(err)
		return
	}
	if err := e.ActivateClient(); err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	w.Header().Set("HX-Redirect", "/dashboard")
}
