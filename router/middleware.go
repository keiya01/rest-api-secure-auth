package router

import (
	"github.com/gorilla/csrf"
	"github.com/keiya01/rest-api-secure-auth/crypto"
	"net/http"
)

func useCSRF() func(http.Handler) http.Handler {
	// TODO: Remove csrf.Secure in production
	return csrf.Protect(crypto.GenerateRandomKey(32), csrf.Secure(false))
}

func (r *Router) middleware() {
	r.Use(useCSRF())
}
