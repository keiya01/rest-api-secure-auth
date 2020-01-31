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
			if !a.AutoLogin(w, r) {
				next.ServeHTTP(w, r)
			}
		})
	})
	authRouter.HandleFunc("/token", a.Token).Methods(http.MethodGet)
	authRouter.HandleFunc("/signup", a.SignUp).Methods(http.MethodPost)
	authRouter.HandleFunc("/{provider}/callback", a.AuthCallback).Methods(http.MethodGet)
	authRouter.HandleFunc("/{provider}", a.ExternalLogin).Methods(http.MethodGet)
	authRouter.HandleFunc("/", a.Login).Methods(http.MethodPost)
}
