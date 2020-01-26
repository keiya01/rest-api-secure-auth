package router

import (
	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func newRouter() *Router {
	return &Router{
		mux.NewRouter(),
	}
}

func (r *Router)newSubrouter(path string) *Router {
	return &Router{
		r.PathPrefix("/api/v1").Subrouter(),
	}
}

func UseRouter() *Router {
	r := newRouter()
	
	r.middleware()

	api := r.newSubrouter("/api/v1")
	api.auth()
	api.user()

	return r
}
