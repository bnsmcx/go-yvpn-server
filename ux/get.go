package ux

import (
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"yvpn_server/db"
	"yvpn_server/do"
)

type pageData struct {
	LoggedIn bool
	Account  *db.Account
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
	id := r.Context().Value("id")
	if id != nil {
		pd.LoggedIn = true
	}

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func RenderSignup(w http.ResponseWriter, r *http.Request) {
	// Parse the templates
	tmpl, err := template.ParseFiles("templates/base.html", "templates/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if user is logged in, add to page data
	var pd pageData
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
	// Get the user account
	a, err := db.GetAccount(r.Context().Value("id").(uuid.UUID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// populate page data
	pd := pageData{
		LoggedIn: true,
		Account:  a,
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
	a, err := db.GetAccount(r.Context().Value("id").(uuid.UUID))
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

func RenderNewCreditNode(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	// Parse the templates
	tmpl, err := template.ParseFiles("templates/base.html", "templates/credit_node.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n, err := db.GetAccount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, n); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
