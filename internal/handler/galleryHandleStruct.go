package handler

type gallery struct {
	ID          int    `json:"id" binding:"required"`
	Photo       string `json:"photo" binding:"required"`
	Description string `json:"description"`
}
