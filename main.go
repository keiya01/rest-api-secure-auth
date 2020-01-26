package main

import (
	"github.com/keiya01/rest-api-secure-auth/database"
	"github.com/keiya01/rest-api-secure-auth/sessions"
	"github.com/keiya01/rest-api-secure-auth/router"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/keiya01/rest-api-secure-auth/auth"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	auth.SetProvider()

	sessionStore := sessions.NewStore()
	sessions.SetSessionStore(sessionStore)

	gothic.Store = sessions.NewStore()

	database.SetDB(database.NewDB())
}

func main() {
	fmt.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", router.UseRouter()))
}
