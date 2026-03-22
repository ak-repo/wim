package middleware

import "net/http"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "authentication middleware not implemented", http.StatusNotImplemented)
	})
}
