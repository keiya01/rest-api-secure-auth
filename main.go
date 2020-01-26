package main

import (
	"github.com/keiya01/rest-api-sample/handler"
	"github.com/keiya01/rest-api-sample/database"
	"github.com/keiya01/rest-api-sample/sessions"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/keiya01/rest-api-sample/auth"
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

func userProfile(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)

	userID, ok := pathParams["userID"];
	if !ok {
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}

	user := database.Get(userID)
	if user == nil {
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"message": "User found"
		"user": %v
	}`, user)))
}

func setHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()

	r.Use(setHeader)

	a := handler.NewAuthHandler()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users/{userID}", userProfile).Methods(http.MethodGet)
	api.HandleFunc("/login/{provider}", a.Login).Methods(http.MethodGet)
	api.HandleFunc("/login/{provider}/callback", a.AuthCallback).Methods(http.MethodGet)

	fmt.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
