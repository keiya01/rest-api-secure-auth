package router

import (
	"github.com/keiya01/rest-api-secure-auth/handler"
	"net/http"
)

func (r *Router) auth() {
	authRouter := r.PathPrefix("/auth").Subrouter()
	a := handler.NewAuthHandler()
	authRouter.HandleFunc("/token", a.Token).Methods(http.MethodGet)
	authRouter.HandleFunc("/signup", a.AutoLogin(a.SignUp)).Methods(http.MethodPost)
	authRouter.HandleFunc("/{provider}/callback", a.AuthCallback).Methods(http.MethodGet)
	authRouter.HandleFunc("/{provider}", a.AutoLogin(a.ExternalLogin)).Methods(http.MethodGet)
	authRouter.HandleFunc("/", a.AutoLogin(a.Login)).Methods(http.MethodPost)
}
