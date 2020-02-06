package sessions

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var CookieOptions = &sessions.Options{
	Path:     "/",
	HttpOnly: true,
	Secure:   false,
	MaxAge:   86400 * 30,
	SameSite: http.SameSiteNoneMode,
}
