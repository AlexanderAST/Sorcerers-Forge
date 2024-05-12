package handler

type profileInput struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
	Contact    string `json:"contact" binding:"required"`
	Photo      string `json:"photo"`
}
