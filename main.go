package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/keiya01/rest-api-sample/auth"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
	"strconv"
)

func login(w http.ResponseWriter, r *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		w.Write([]byte(fmt.Sprintf(`{
			"user": %v
		}`, gothUser)))
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func authCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		w.Write([]byte(`{"message": "Login failure ..."}`))
		return
	}
	w.Write([]byte(fmt.Sprintf(`{
		"message": "Login Success!",
		"user": %v
	}`, user)))
}

func userProfile(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)

	var userID int
	var err error

	if val, ok := pathParams["userID"]; ok {
		userID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "user id param need to be number"}`))
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"user": {
			"id": %d,
			"name": "user%d"
		}
	}`, userID, userID)))
}

func json(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	auth.SetProvider()

	r := mux.NewRouter()

	r.Use(json)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users/{userID}", userProfile).Methods(http.MethodGet)
	api.HandleFunc("/login/{provider}", login).Methods(http.MethodGet)
	api.HandleFunc("/login/{provider}/callback", authCallback).Methods(http.MethodGet)

	fmt.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
