package model

type Reviews struct {
	ID        int    `json:"id"`
	ProductId int    `json:"product_id"`
	UserID    int    `json:"user_id"`
	Stars     int    `json:"stars" binding:"required"`
	Message   string `json:"message"`
}
