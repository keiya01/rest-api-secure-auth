package response

import (
	"github.com/gorilla/csrf"
	"net/http"
)

func UseCSRFToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
}

func UseJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SetHeader(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	UseJSON(w)
}
