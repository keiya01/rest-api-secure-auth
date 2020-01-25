package sessions

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func NewStore() *sessions.CookieStore {
	return sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
}
