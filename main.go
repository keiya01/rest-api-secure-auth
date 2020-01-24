package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

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
	r := mux.NewRouter()
	r.Use(json)
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users/{userID}", userProfile).Methods(http.MethodGet)
	fmt.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
