package handler

type category struct {
	Name string `json:"name" binding:"required"`
}

type reqWithIDProduct struct {
	ID int `json:"id" binding:"required"`
}

type catalog struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Price        int     `json:"price" binding:"required"`
	ReviewsMid   float64 `json:"reviews_mid"`
	ReviewsCount int     `json:"reviews_count"`
	Quantity     int     `json:"quantity"`
	WorkTime     int     `json:"work_time"`
	Photo        string  `json:"photo"`
	CategoryID   int     `json:"category_id" binding:"required"`
	IsActive     bool    `json:"is_active"`
}

type cartItems struct {
	ProductID int    `json:"product_id" binding:"required"`
	Count     int    `json:"count" binding:"required"`
	Photo     string `json:"photo"`
}

type favoriteItems struct {
	ProductID int    `json:"product_id" binding:"required"`
	Photo     string `json:"photo"`
}

type updateCatalog struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	WorkTime    int    `json:"work_time"`
	Photo       string `json:"photo"`
	CategoryID  int    `json:"category_id" binding:"required"`
	IsActive    bool   `json:"is_active"`
}

type reqLikes struct {
	ID        int `json:"id"`
	ProductID int `json:"product_id"`
}
