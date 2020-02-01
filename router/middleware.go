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
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Accept, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
		w.Header().Set("Access-Control-Max-Age", "86400") // Time to cache preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// Use X-Frame-Options as middleware
func common(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Available values are DENY, SAMEORIGIN, ALLOW-FROM
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Xss-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

func (r *Router) middleware() {
	r.Use(
		cors,
		common,
		csrf.Protect(
			crypto.GenerateRandomKey(32),
			csrf.Secure(false), // TODO: Remove csrf.Secure in production
			csrf.TrustedOrigins([]string{"localhost:3000"}),
		),
	)
}
