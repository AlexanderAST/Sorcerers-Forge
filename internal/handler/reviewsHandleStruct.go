package handler

type reviews struct {
	ID        int    `json:"id" binding:"required" `
	ProductID int    `json:"product_id"`
	UserID    int    `json:"user_id" binding:"required"`
	Stars     int    `json:"stars" binding:"required"`
	Message   string `json:"message"`
}
