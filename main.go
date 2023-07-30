package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"yvpn_server/auth"
	"yvpn_server/db"
	"yvpn_server/ux"
)

func main() {
	r := chi.NewRouter()
	err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Get("/", ux.RenderLanding)
		r.Get("/signup", ux.RenderSignup)
		r.Get("/login", ux.RenderLogin)
		r.Get("/logout", auth.HandleLogout)
	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(auth.UserSession)
		r.Get("/manage", ux.RenderManage)
	})

	log.Fatalln(http.ListenAndServe(":8000", r))
}
