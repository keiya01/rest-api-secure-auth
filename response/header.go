package response

import (
	"github.com/gorilla/csrf"
	"net/http"
)

func UseJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SetAuthAPIHeader(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.WriteHeader(status)
	UseJSON(w)
}

func SetOpenAPIHeader(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	UseJSON(w)
}
