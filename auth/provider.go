package auth

import (
	"github.com/keiya01/rest-api-secure-auth/model"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"net/http"
	"os"
)

func SetProvider() {
	goth.UseProviders(twitter.New(os.Getenv("TWITTER_CLIENT_KEY"), os.Getenv("TWITTER_SECRET_KEY"), "http://localhost:8080/api/v1/login/twitter/callback"))
}

func AuthProvider(w http.ResponseWriter, r *http.Request) (model.User, bool) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return model.User{}, false
	}

	return model.NewUser(gothUser.UserID,gothUser.Name,gothUser.Description), true
}
