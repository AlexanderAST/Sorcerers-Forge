package model

type Message struct {
	UserEmail string `json:"user_email" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Message   string `json:"message" binding:"required"`
}
