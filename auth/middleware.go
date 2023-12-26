package auth

import (
	"context"
	"log"
	"net/http"
)

// MandateSession requires a user has a valid session token
func MandateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Println("invalid or missing session cookie: ", err)
			redirectToLogin(w, r)
			return
		}
		id, ok := getSession(cookie.Value)
		if !ok {
			log.Println("invalid session cookie: ", cookie.Value)
			redirectToLogin(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "id", id)
		ctx = context.WithValue(ctx, "session_id", cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Set("HX-Redirect", "/login")
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

// CheckSession checks to see if a user is logged in without forcing a redirect to /login
func CheckSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("session_id")
		if err == nil {
			id, ok := getSession(cookie.Value)
			if ok {
				ctx = context.WithValue(ctx, "id", id)
				ctx = context.WithValue(ctx, "session_id", cookie.Value)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
