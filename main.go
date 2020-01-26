package main

import (
	"github.com/keiya01/rest-api-sample/database"
	"github.com/keiya01/rest-api-sample/sessions"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/keiya01/rest-api-sample/auth"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
)

var (
	SESSION_STORE_NAME = "cookie-store"
	db = database.NewDB()
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
}

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

type loginResponse struct {
	Message string `json:"message"`
	User User `json:"user"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var err error

	session, _ := sessions.Get(r, SESSION_STORE_NAME)
	if userID, ok := session.Values["userID"].(string); ok {
			if user, ok := db.Get(userID).(User); ok {
				userJSON, _ := json.Marshal(loginResponse {
					Message: "Auto login success",
					User: user,
				})
				w.Write(userJSON)
				return
			}
	}

	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}

	userJSON, err := json.Marshal(loginResponse {
		Message: "Login success",
		User: User {
			ID: gothUser.UserID,
			Name: gothUser.Name,
			Description: gothUser.Description,
		},
	})
	if err != nil {
		w.Write([]byte(`{message: "Could not be received response data"}`))
		return
	}

	w.Write(userJSON)
}

func authCallback(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		w.Write([]byte(`{"message": "Login failure ..."}`))
		return
	}

	session, _ := sessions.Get(r, SESSION_STORE_NAME)

	user := User {
		ID: gothUser.UserID,
		Name: gothUser.Name,
		Description: gothUser.Description,
	}

	session.Values["userID"] = user.ID
	session.Options = sessions.CookieOptions
	sessions.Save(r, w, session)

	db.Insert(user.ID, user)

	vars := mux.Vars(r)
	http.Redirect(
		w, 
		r, 
		fmt.Sprintf("/api/v1/login/%s", vars["provider"]), 
		http.StatusTemporaryRedirect,
	)
}

func userProfile(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)

	userID, ok := pathParams["userID"];
	if !ok {
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}

	user := db.Get(userID)
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

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users/{userID}", userProfile).Methods(http.MethodGet)
	api.HandleFunc("/login/{provider}", login).Methods(http.MethodGet)
	api.HandleFunc("/login/{provider}/callback", authCallback).Methods(http.MethodGet)

	fmt.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
