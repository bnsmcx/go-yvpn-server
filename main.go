package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"strings"
	"yvpn_server/auth"
	"yvpn_server/db"
	"yvpn_server/ux"
)

func main() {
	err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// Create the router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.DefaultLogger)

	// Serve static files from the "static" directory
	fileServer(r, "/static", http.Dir("./static"))

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Get("/", ux.RenderLanding)
		r.Get("/signup", ux.RenderSignup)
		r.Get("/login", ux.RenderLogin)
		r.Get("/logout", auth.HandleLogout)
		r.Post("/signup", ux.SignupHandler)
	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(auth.UserSession)
		r.Get("/manage", ux.RenderManage)
	})

	log.Fatalln(http.ListenAndServe(":8000", r))
}

func fileServer(r *chi.Mux, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit any URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
