package auth

import (
	"context"
	"log"
	"net/http"
)

// UserSession validates a user's session token and permissions
func UserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Println("missing session_id")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		id, ok := getSession(cookie.Value)
		if !ok {
			log.Println("invalid session_id")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		ctx := context.WithValue(r.Context(), "id", id)
		ctx = context.WithValue(ctx, "session_id", cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
