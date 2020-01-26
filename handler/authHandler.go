package handler

import (
	"github.com/gorilla/mux"
	"fmt"
	"github.com/markbates/goth/gothic"
	"encoding/json"
	"github.com/keiya01/rest-api-sample/sessions"
	"net/http"
	"github.com/keiya01/rest-api-sample/database"
	"github.com/keiya01/rest-api-sample/model"
)

type AuthHandler struct {}

type loginResponse struct {
	Message string `json:"message"`
	User model.User `json:"user"`
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (a *AuthHandler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		w.Write([]byte(`{"message": "Login failure ..."}`))
		return
	}

	session, _ := sessions.Get(r, sessions.SESSION_STORE_NAME)

	user := model.User {
		ID: gothUser.UserID,
		Name: gothUser.Name,
		Description: gothUser.Description,
	}

	session.Values["userID"] = user.ID
	session.Options = sessions.CookieOptions
	sessions.Save(r, w, session)

	database.Insert(user.ID, user)

	vars := mux.Vars(r)
	http.Redirect(
		w, 
		r, 
		fmt.Sprintf("/api/v1/login/%s", vars["provider"]), 
		http.StatusTemporaryRedirect,
	)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var err error

	session, _ := sessions.Get(r, sessions.SESSION_STORE_NAME)
	if userID, ok := session.Values["userID"].(string); ok {
			if user, ok := database.Get(userID).(model.User); ok {
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
		User: model.User {
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
