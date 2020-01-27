package sessions

import (
	"github.com/gorilla/sessions"
	"github.com/keiya01/rest-api-secure-auth/crypto"
	"net/http"
)

var SESSION_STORE_NAME = "cookie-store"

var store = &sessions.CookieStore{}

func NewStore() *sessions.CookieStore {
	return sessions.NewCookieStore(crypto.GenerateRandomKey(64), crypto.GenerateRandomKey(32))
}

func SetSessionStore(sessionStore *sessions.CookieStore) {
	store = sessionStore
}

func Get(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}

func Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return store.Save(r, w, session)
}
