package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Received GET method"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Received POST method"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "Received PUT method"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Received DELETE method"}`))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "Received method not found"}`))
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
	api.HandleFunc("/", get).Methods(http.MethodGet)
	api.HandleFunc("/", post).Methods(http.MethodPost)
	api.HandleFunc("/", put).Methods(http.MethodPut)
	api.HandleFunc("/", delete).Methods(http.MethodDelete)
	api.HandleFunc("/", notFound)
	fmt.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
