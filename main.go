package main

import (
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

type server struct{}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Received GET method"}`))
	case "POST":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Received POST method"}`))
	case "PUT":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Received PUT method"}`))
	case "DELETE":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Received DELETE method"}`))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Received method not found"}`))
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	fmt.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
