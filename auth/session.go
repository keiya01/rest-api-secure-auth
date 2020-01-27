package auth

import (
	"github.com/keiya01/rest-api-secure-auth/database"
	"github.com/keiya01/rest-api-secure-auth/model"
	"github.com/keiya01/rest-api-secure-auth/sessions"
	"net/http"
)

func checkSession(key string, r *http.Request) (string, bool) {
	session, _ := sessions.Get(r, sessions.SESSION_STORE_NAME)
	v, ok := session.Values[key]
	if !ok {
		return "", false
	}

	val, ok := v.(string)
	if !ok {
		return "", false
	}

	return val, true
}

func IsLogin(sessionKey string, r *http.Request) (model.User, bool) {
	userID, ok := checkSession(sessionKey, r)
	if !ok {
		return model.User{}, false
	}

	user, ok := database.Get(userID).(model.User)
	if !ok {
		return model.User{}, false
	}

	return user, true
}
