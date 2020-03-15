package graph

import (
	"context"
	"net/http"
)

// Middleware is an example usage of a simple middlware.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "Message", "Hello World")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// CorsMiddleware is a help function get around CORS
// TODO do this better with authentication
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
