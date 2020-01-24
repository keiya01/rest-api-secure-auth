package auth

import (
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/twitter"
	"log"
	"os"
)

func SetProvider() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	goth.UseProviders(twitter.New(os.Getenv("TWITTER_CLIENT_KEY"), os.Getenv("TWITTER_SECRET_KEY"), "http://localhost:8080/api/v1/login/twitter/callback"))
}
