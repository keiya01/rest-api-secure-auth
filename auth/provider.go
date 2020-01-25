package auth

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/twitter"
	"os"
)

func SetProvider() {
	goth.UseProviders(twitter.New(os.Getenv("TWITTER_CLIENT_KEY"), os.Getenv("TWITTER_SECRET_KEY"), "http://localhost:8080/api/v1/login/twitter/callback"))
}
