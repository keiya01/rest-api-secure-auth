package router

import (
	"github.com/gorilla/csrf"
	"github.com/keiya01/rest-api-secure-auth/crypto"
	"net/http"
)

func useCSRF() func(http.Handler) http.Handler {
	return csrf.Protect(crypto.GenerateRandomKey(32))
}

func (r *Router) middleware() {
	r.Use(useCSRF())
}
