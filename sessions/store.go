package sessions

import (
	"net/http"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var store = &sessions.CookieStore{}

func NewStore() *sessions.CookieStore {
	return sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
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
