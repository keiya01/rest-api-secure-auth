package router

import (
	"github.com/gorilla/csrf"
	"github.com/keiya01/rest-api-secure-auth/crypto"
	"net/http"
)

func useCSRF() func(http.Handler) http.Handler {
	return csrf.Protect(crypto.GenerateRandomKey(32))
}

func useJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (r *Router) middleware() {
	r.Use(useJSON)
	r.Use(useCSRF())
}
