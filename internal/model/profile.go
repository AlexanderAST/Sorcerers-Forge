package model

type Profile struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
	Contact    string `json:"contact" binding:"required"`
	Photo      string `json:"photo"`
}
