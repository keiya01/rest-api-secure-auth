package router

import (
	"net/http"
	"github.com/keiya01/rest-api-secure-auth/handler"
)

func (r *Router) auth() {
	a := handler.NewAuthHandler()
	r.HandleFunc("/login/{provider}", a.Login).Methods(http.MethodGet)
	r.HandleFunc("/login/{provider}/callback", a.AuthCallback).Methods(http.MethodGet)
}
