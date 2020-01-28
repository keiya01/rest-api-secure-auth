package router

import (
	"github.com/keiya01/rest-api-secure-auth/handler"
	"net/http"
)

func (r *Router) auth() {
	authRouter := r.PathPrefix("/auth").Subrouter()
	a := handler.NewAuthHandler()
	authRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !a.HandleAutoLogin(w, r) {
				next.ServeHTTP(w, r)
			}
		})
	})
	authRouter.HandleFunc("/{provider}", a.HandleExternalLogin).Methods(http.MethodGet)
	authRouter.HandleFunc("/", a.HandleLogin).Methods(http.MethodPost)
	authRouter.HandleFunc("/{provider}/callback", a.AuthCallback).Methods(http.MethodGet)
}
