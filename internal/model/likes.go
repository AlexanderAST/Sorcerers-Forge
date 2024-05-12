package model

type Like struct {
	ID        int `json:"id"`
	UserId    int `json:"user_id"`
	ProductID int `json:"product_id"`
}
