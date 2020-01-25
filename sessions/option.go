package sessions

import (
	"github.com/gorilla/sessions"
)

var CookieOptions = &sessions.Options{
	Path: "/",
	HttpOnly: true,
	Secure: false,
	MaxAge: 86400 * 30,
}
