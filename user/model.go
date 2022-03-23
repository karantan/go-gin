package user

import (
	"errors"
)

//--
// Data model objects and persistence mocks:
//--

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

var users = []*User{
	{ID: 1, Email: "admin", Password: "secret"},
	{ID: 2, Email: "foo", Password: "secret"},
}

func DBGetUser(id int) (*User, error) {
	for _, u := range users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found.")
}

func DBGetUserByEmail(email string) (*User, error) {
	for _, u := range users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found.")
}

func DBUpdateUser(id int, user *User) (*User, error) {
	for i, u := range users {
		if u.ID == id {
			users[i] = user
			return user, nil
		}
	}
	return nil, errors.New("user not found.")
}
