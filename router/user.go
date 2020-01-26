package router

import (
	"github.com/keiya01/rest-api-sample/handler"
	"net/http"
)

func (r *Router) user() {
	u := handler.NewUserHandler()
	r.HandleFunc("/users/{userID}", u.Profile).Methods(http.MethodGet)
}
