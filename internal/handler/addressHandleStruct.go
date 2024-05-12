package handler

type address struct {
	ID     int    `json:"id"`
	Name   string `json:"name" binding:"required"`
	Latlng string `json:"latlng" binding:"required"`
}
