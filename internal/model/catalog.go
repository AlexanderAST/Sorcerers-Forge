package model

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

type Products struct {
	ID           int     `json:"id"`
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

type CartItem struct {
	ID        int    `json:"id"`
	UserId    int    `json:"user_id" binding:"required"`
	ProductID int    `json:"product_id" binding:"required"`
	Count     int    `json:"count" binding:"required"`
	Photo     string `json:"photo"`
}

type FavoriteItem struct {
	ID        int    `json:"id"`
	UserId    int    `json:"user_id" binding:"required"`
	ProductID int    `json:"product_id" binding:"required"`
	Photo     string `json:"photo"`
}

type Cart struct {
	ID          int    `json:"id"`
	UserId      int    `json:"user_id" binding:"required"`
	ProductID   int    `json:"product_id" binding:"required"`
	UserEmail   string `json:"user_email"`
	ProductName string `json:"product_name"`
	Count       int    `json:"count" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Photo       string `json:"photo"`
}

type Favorite struct {
	ID          int    `json:"id"`
	UserId      int    `json:"user_id" binding:"required"`
	ProductID   int    `json:"product_id" binding:"required"`
	UserEmail   string `json:"user_email"`
	ProductName string `json:"product_name"`
	Price       int    `json:"price" binding:"required"`
	Photo       string `json:"photo"`
}

type ProductWithFlags struct {
	*Products
	IsFavorite bool
	IsCart     bool
}

type Orders struct {
	ID           int `json:"id"`
	UserID       int `json:"user_id"`
	Products     []ProductCartInfo
	Summ         int   `json:"summ"`
	ProductID    []int `json:"product_id"`
	ProductCount []int `json:"product_count"`
}

type ProductCartInfo struct {
	ProductId int `json:"product_id"`
	Count     int `json:"count"`
}
