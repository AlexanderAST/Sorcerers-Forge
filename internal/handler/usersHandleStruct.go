package handler

import "github.com/dgrijalva/jwt-go"

type reqWithEmailInput struct {
	Email string `json:"email" binding:"required"`
}

type request struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type resetPassword struct {
	Email     string `json:"email" binding:"required"`
	EmailCode string `json:"emailCode" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}
