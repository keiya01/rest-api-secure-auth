package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/keiya01/rest-api-secure-auth/auth"
	"github.com/keiya01/rest-api-secure-auth/crypto"
	"github.com/keiya01/rest-api-secure-auth/database"
	"github.com/keiya01/rest-api-secure-auth/model"
	"github.com/keiya01/rest-api-secure-auth/response"
	"github.com/keiya01/rest-api-secure-auth/sessions"
	"github.com/markbates/goth/gothic"
	"io/ioutil"
	"net/http"
)

type AuthHandler struct{}

type loginResponse struct {
	Message  string     `json:"message"`
	User     model.User `json:"user"`
	Provider string     `json:"provider"`
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (a *AuthHandler) HandleAutoLogin(w http.ResponseWriter, r *http.Request) bool {
	if currentUser, ok := auth.IsLogin("userID", r); ok {
		userJSON, _ := json.Marshal(loginResponse{
			Message:  "Auto login success",
			User:     currentUser,
			Provider: mux.Vars(r)["provider"],
		})
		response.SetAuthAPIHeader(w, r, http.StatusOK)
		w.Write(userJSON)
		return true
	}
	return false
}

func (a *AuthHandler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		response.UseCSRFToken(w, r)
		http.Error(w, "Login failure", http.StatusInternalServerError)
		return
	}

	session, _ := sessions.Get(r, sessions.SESSION_STORE_NAME)

	user := model.NewUser(gothUser.UserID, gothUser.Name, gothUser.Description, "")

	session.Values["userID"] = user.ID
	session.Options = sessions.CookieOptions
	sessions.Save(r, w, session)

	database.Insert(user.ID, user)

	http.Redirect(
		w,
		r,
		// TODO: FrontEndのpathに変更する
		fmt.Sprintf("/api/v1/users/%s", user.ID),
		http.StatusTemporaryRedirect,
	)
}

func (a *AuthHandler) HandleExternalLogin(w http.ResponseWriter, r *http.Request) {
	if _, err := gothic.GetProviderName(r); err != nil {
		response.UseCSRFToken(w, r)
		http.Error(w, "The request has Invalid parameter", http.StatusBadRequest)
		return
	}

	user, ok := auth.AuthProvider(w, r)
	if !ok {
		gothic.BeginAuthHandler(w, r)
		return
	}

	userJSON, err := json.Marshal(loginResponse{
		Message:  "Login success",
		User:     user,
		Provider: mux.Vars(r)["provider"],
	})

	if err != nil {
		response.UseCSRFToken(w, r)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	response.SetAuthAPIHeader(w, r, http.StatusOK)
	w.Write(userJSON)
}

type LoginUser struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

func (a *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response.UseCSRFToken(w, r)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user LoginUser
	err = json.Unmarshal(b, &user)
	if err != nil {
		response.UseCSRFToken(w, r)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(user)

	errors := []map[string]string{}

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMap := map[string]string{}
			field := err.Field()
			switch field {
			case "Username":
				errorMap = map[string]string{
					"field":     "username",
					"errorType": err.Type().String(),
				}
			case "Email":
				errorMap = map[string]string{
					"field":     "email",
					"errorType": err.Type().String(),
				}
			case "Password":
				errorMap = map[string]string{
					"field":     "username",
					"errorType": err.Type().String(),
				}
			}
			errors = append(errors, errorMap)
		}
		response.SetAuthAPIHeader(w, r, http.StatusBadRequest)
		errorRes := map[string]interface{}{
			"errors": errors,
		}
		errorResJSON, _ := json.Marshal(errorRes)
		w.Write(errorResJSON)
		return
	}

	// TODO: Add email authentication
	session, _ := sessions.Get(r, sessions.SESSION_STORE_NAME)

	resUser := model.NewUser(string(crypto.GenerateRandomKey(32)), user.Username, "", user.Email)

	session.Values["userID"] = resUser.ID
	session.Options = sessions.CookieOptions
	sessions.Save(r, w, session)

	database.Insert(resUser.ID, user)

	res := loginResponse{
		Message:  "Login Success",
		User:     resUser,
		Provider: "",
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		response.UseCSRFToken(w, r)
		http.Error(w, "Failed JSON encoding of response", http.StatusInternalServerError)
		return
	}
	w.Write(resJSON)
}
