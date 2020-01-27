package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/keiya01/rest-api-secure-auth/database"
	"net/http"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)

	userID, ok := pathParams["userID"]
	if !ok {
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}

	user := database.Get(userID)
	if user == nil {
		w.Write([]byte(`{"message": "User not found"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"message": "User found"
		"user": %v
	}`, user)))
}
