package model

type Gallery struct {
	ID          int    `json:"id"`
	Photo       string `json:"catalog" binding:"required"`
	Description string `json:"description"`
}
