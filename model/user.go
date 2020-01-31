package model

import (
	"github.com/keiya01/rest-api-secure-auth/database"
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func NewUser(id string, name string, description string, email string, password string) User {
	return User{
		ID:          id,
		Name:        name,
		Description: description,
		Email:       email,
		Password:    password,
	}
}

func (u User) FindByEmail() User {
	db := database.GetDB()
	for _, item := range db.Values {
		user, ok := item.(User)
		if !ok {
			return User{}
		}
		if user.Email == u.Email {
			return user
		}
	}
	return User{}
}
