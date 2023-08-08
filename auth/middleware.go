package auth

import (
	"context"
	"net/http"
)

// UserSession validates a user's session token and permissions
func UserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "missing session_id", http.StatusUnauthorized)
		}
		id, ok := getSession(cookie.Value)
		if !ok {
			http.Error(w, "invalid session_id", http.StatusUnauthorized)
		}
		ctx := context.WithValue(r.Context(), "id", id)
		ctx = context.WithValue(ctx, "session_id", cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
