package router

import (
	"net/http"
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

func setHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (r *Router) middleware() {
	r.Use(setHeader)
}

func UseRouter() *Router {
	r := newRouter()
	
	r.middleware()

	api := r.newSubrouter("/api/v1")
	api.auth()
	api.user()

	return r
}
