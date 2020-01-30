package router

import (
	"github.com/gorilla/csrf"
	"github.com/keiya01/rest-api-secure-auth/crypto"
	"net/http"
)

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Accept, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400") // Time to cache preflight request
		next.ServeHTTP(w, r)
	})
}

func (r *Router) middleware() {
	r.Use(cors)
	r.Use(csrf.Protect(crypto.GenerateRandomKey(32), csrf.Secure(false))) // TODO: Remove csrf.Secure in production
}
