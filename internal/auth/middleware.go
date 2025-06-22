package auth

import (
	"net/http"
	"os"
)

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := r.URL.Query().Get("admin")
		expected := os.Getenv("DASHBOARD_PASSWORD")
		if expected == "" {
			expected = "admin"
		}

		if password != expected {
			http.Error(w, "Unauthorized access to dashboard", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
