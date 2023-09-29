package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"log"
	"net/http"
	"strings"
	"time"
	"yvpn_server/auth"
	"yvpn_server/db"
	"yvpn_server/do"
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
	r.Use(httprate.LimitByIP(60, 1*time.Minute))

	// Serve static files from the "static" directory
	fileServer(r, "/static", http.Dir("./static"))

	// Redirect all 404's to home page
	r.NotFound(ux.RenderLanding)

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Use(auth.CheckSession)

		r.Get("/", ux.RenderLanding)
		r.Get("/signup", ux.RenderSignup)
		r.Get("/activate", ux.RenderActivation)
		r.Get("/login", ux.RenderLogin)

		r.Post("/activate", ux.ActivationHandler)
		r.Post("/signup", ux.SignupHandler)
		//r.Post("/login", auth.Login)
	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(auth.MandateSession)

		r.Get("/dashboard", ux.RenderDashboard)
		r.Get("/logout", auth.HandleLogout)
		r.Get("/token/add", ux.RenderAddToken)
		r.Get("/endpoints/add", ux.RenderAddEndpoint)
		r.Get("/endpoints/{endpoint}/{client}", db.GetClientConfigFile)

		r.Delete("/endpoints/{id}", do.HandleDeleteEndpoint)

		r.Post("/endpoints/add", do.HandleCreateEndpoint)
		r.Post("/token/add", do.AddToken)
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
