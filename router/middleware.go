package router

import (
	"net/http"
)

func setHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (r *Router) middleware() {
	r.Use(setHeader)
}
