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
	"github.com/keiya01/rest-api-secure-auth/validation"
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

func (a *AuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	response.UseCSRFToken(w, r)
	w.WriteHeader(http.StatusOK)
}

func (a *AuthHandler) AutoLogin(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if currentUser, ok := auth.IsLogin("userID", r); ok {
			userJSON, _ := json.Marshal(loginResponse{
				Message:  "Auto login success",
				User:     model.NewUser(currentUser.ID, currentUser.Name, currentUser.Description, currentUser.Email, ""),
				Provider: mux.Vars(r)["provider"],
			})
			response.SetHeader(w, r, http.StatusOK)
			w.Write(userJSON)
			return
		}
		f(w, r)
	})
}

func (a *AuthHandler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		response.Error(w, "Login failure", http.StatusInternalServerError)
		return
	}

	session, err := sessions.Get(r, sessions.SESSION_STORE_NAME)
	if err != nil {
		sessions.SetCookieStore(sessions.NewStore())
	}

	user := model.NewUser(gothUser.UserID, gothUser.Name, gothUser.Description, "", "")

	session.Values["userID"] = user.ID
	session.Options = sessions.CookieOptions
	sessions.Save(r, w, session)

	database.Insert(user.ID, user)

	http.Redirect(
		w,
		r,
		// TODO: Change URL for the front end
		fmt.Sprintf("/api/v1/users/%s", user.ID),
		http.StatusTemporaryRedirect,
	)
}

func (a *AuthHandler) ExternalLogin(w http.ResponseWriter, r *http.Request) {
	if _, err := gothic.GetProviderName(r); err != nil {
		response.Error(w, "The request has Invalid parameter", http.StatusBadRequest)
		return
	}

	user, ok := auth.AuthProvider(w, r)
	if !ok {
		gothic.BeginAuthHandler(w, r)
		return
	}

	userJSON, err := json.Marshal(loginResponse{
		Message:  "Login success",
		User:     model.NewUser(user.ID, user.Name, user.Description, user.Email, ""),
		Provider: mux.Vars(r)["provider"],
	})

	if err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	response.SetHeader(w, r, http.StatusOK)
	w.Write(userJSON)
}

type SignUpUser struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user SignUpUser
	err = json.Unmarshal(b, &user)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(user)

	if err != nil {
		errorRes := map[string]interface{}{
			"errors": validation.Extract(err.(validator.ValidationErrors), []string{"Username", "Email", "Password"}),
		}
		errorResJSON, _ := json.Marshal(errorRes)
		response.SetHeader(w, r, http.StatusBadRequest)
		w.Write(errorResJSON)
		return
	}

	session, err := sessions.Get(r, sessions.SESSION_STORE_NAME)
	if err != nil {
		sessions.SetCookieStore(sessions.NewStore())
	}

	resUser := model.NewUser(string(crypto.GenerateRandomKey(32)), user.Username, "", user.Email, user.Password)

	session.Values["userID"] = resUser.ID
	session.Options = sessions.CookieOptions
	sessions.Save(r, w, session)

	database.Insert(resUser.ID, resUser)

	res := loginResponse{
		Message:  "Login Success",
		User:     model.NewUser(resUser.ID, resUser.Name, resUser.Description, resUser.Email, ""),
		Provider: "",
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		response.Error(w, "Failed JSON encoding of response", http.StatusInternalServerError)
		return
	}
	response.SetHeader(w, r, http.StatusOK)
	w.Write(resJSON)
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var loginUser LoginUser
	err = json.Unmarshal(b, &loginUser)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(loginUser)

	if err != nil {
		errorRes := map[string]interface{}{
			"errors": validation.Extract(err.(validator.ValidationErrors), []string{"Email", "Password"}),
		}
		errorResJSON, _ := json.Marshal(errorRes)
		response.SetHeader(w, r, http.StatusBadRequest)
		w.Write(errorResJSON)
		return
	}

	var (
		currentUser model.User
	)

	user := model.NewUser("", "", "", loginUser.Email, loginUser.Password)
	currentUser = user.FindByEmail()

	if currentUser.Password != user.Password {
		response.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	session, err := sessions.Get(r, sessions.SESSION_STORE_NAME)
	if err != nil {
		sessions.SetCookieStore(sessions.NewStore())
	}

	session.Values["userID"] = currentUser.ID
	session.Options = sessions.CookieOptions
	sessions.Save(r, w, session)

	res := loginResponse{
		Message:  "Login Success",
		User:     model.NewUser(currentUser.ID, currentUser.Name, currentUser.Description, currentUser.Email, ""),
		Provider: "",
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		response.Error(w, "Failed JSON encoding of response", http.StatusInternalServerError)
		return
	}
	response.SetHeader(w, r, http.StatusOK)
	w.Write(resJSON)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.IsLogin("userID", r); !ok {
		response.Error(w, "You Don't login", http.StatusBadRequest)
		return
	}

	session, err := sessions.Get(r, sessions.SESSION_STORE_NAME)
	if err != nil {
		response.Error(w, "Please checking if you login", http.StatusBadRequest)
		return
	}

	store := sessions.GetCookieStore()
	store.MaxAge(-1)
	err = store.Save(r, w, session)
	if err != nil {
		response.Error(w, "Failed logout", http.StatusInternalServerError)
		return
	}

	response.SetHeader(w, r, http.StatusOK)
	w.Write([]byte(`{"message": "Success Logout"}`))
}
