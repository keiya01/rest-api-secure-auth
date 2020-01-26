package handler

import (
	"fmt"
	"github.com/keiya01/rest-api-sample/database"
	"net/http"
	"github.com/gorilla/mux"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
return &UserHandler{}
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
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
