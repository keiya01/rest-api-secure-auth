package model

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Email       string `json:"email"`
}

func NewUser(id string, name string, description string, email string) User {
	return User{
		ID:          id,
		Name:        name,
		Description: description,
		Email:       email,
	}
}
