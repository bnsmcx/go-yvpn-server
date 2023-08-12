package ux

import (
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"yvpn_server/db"
)

func RenderLanding(w http.ResponseWriter, r *http.Request) {
	// Parse the templates
	tmpl, err := template.ParseFiles("templates/base.html", "templates/landing.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, nil); err != nil {
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

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, nil); err != nil {
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

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, nil); err != nil {
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

	// Parse the templates
	tmpl, err := template.ParseFiles("templates/base.html", "templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the "layout" template and send it to the ResponseWriter.
	if err := tmpl.Execute(w, a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
