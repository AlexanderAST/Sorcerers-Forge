package model

import "testing"

func TestUser(t *testing.T) *User {

	return &User{
		Email:    "email@example.org",
		Password: "qwerty",
	}
}
