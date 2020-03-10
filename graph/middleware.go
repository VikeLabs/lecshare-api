package graph

import (
	"context"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "Message", "Hello World")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
