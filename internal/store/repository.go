package store

import (
	"Diploma/internal/model"
)

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	Find(int) (*model.User, error)
	SendResetCode(email string) (string, error)
	ResetPassword(email, emailCode, password string) (string, error)
	DeleteUser(id int) error
}

type CatalogRepository interface {
	CreateCategory(c *model.Category) error
	GetAllCategories() ([]*model.Category, error)
	DeleteCategory(id int) error
	CreateProduct(products *model.Products) error
	FindByID(id int) (*model.Products, error)
	DeleteProduct(id int) error
	GetAllProducts() ([]*model.Products, error)
	GetAllProductsCategories(categoryID int) ([]*model.Products, error)
	AddToCartProduct(cart *model.CartItem) error
	GetAllCartUser(id int) ([]*model.Cart, error)
	AddToFavoriteProduct(cart *model.FavoriteItem) error
	GetAllFavoriteUser(id int) ([]*model.Favorite, error)
	DeleteCart(productid, userID int) error
	DeleteFavorite(productid, userID int) error
	UpdateProduct(c *model.Products) (string, error)
	GetAllProductsSortedByPriceAsc(offset, limit int) ([]*model.Products, error)
	GetAllProductsSortedByPriceDesc(offset, limit int) ([]*model.Products, error)
	GetAllProductsPaginated(offset, limit int) ([]*model.Products, error)
	AddLikeProduct(l *model.Like) error
	CheckLikes(user_id, product_id int) (bool, error)
	GetProductLikes() ([]*model.Like, error)
	//GetAllProductsSortedByLikesAsc(offset, limit int) ([]*model.Products, error)
	GetAllProductsSortedByLikesDesc(offset, limit int) ([]*model.Products, error)
	GetProductsBetweenPrice(startPrice, secondPrice int, offset, limit int) ([]*model.Products, error)
	GetAllIsActive(offset, limit int, active string) ([]*model.Products, error)
	FilterByType(offset, limit int, category string) ([]*model.Products, error)
	GetFilteredProductsSortedByLikesDesc(offset, limit int, filters map[string]interface{}) ([]*model.Products, error)
	GetAllProductsSortedByLikesDesc1(offset, limit int, priceFilter, activeFilter, typeFilter string) ([]*model.Products, error)
	CountFilteredProducts(priceFilter, activeFilter, typeFilter string) (int, error)
	GetUserFavoriteProducts(userID int) ([]*model.Products, error)
	GetUserCartProducts(userID int) ([]*model.Products, error)
	IsProductInCart(userID, productID int) (bool, error)
	IsProductInFavorites(userID, productID int) (bool, error)
	GetAllProductsWithFlags(userID int) ([]*model.ProductWithFlags, error)
	CreateOrder(cartItems []*model.Cart) (int, error)
	GetUserOrderHistory(userID int) ([]*model.Orders, error)
	DeleteAllFromCart(userID int) error
	GetUserOrderHistoryByReviews(userID, productID int) (bool, error)
}

type ProfileRepository interface {
	CreateProfile(p *model.Profile) error
	UpdateProfile(p *model.Profile) (string, error)
	DeleteProfile(id int) error
	FindByID(id int) (*model.Profile, error)
}

type ReviewsRepository interface {
	CreateReviews(p *model.Reviews) error
	GetAllReviews() ([]*model.Reviews, error)
	UpdateReviews(p *model.Reviews) (string, error)
	DeleteReview(productID, userID int) error
	GetAllReviewsFromProduct(productId int) ([]*model.Reviews, error)
}

type GalleryRepository interface {
	CreateGallery(g *model.Gallery) error
	DeleteGallery(id int) error
	GetAllGallery() ([]*model.Gallery, error)
	UpdateGallery(g *model.Gallery) (string, error)
}

type AdminRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	Find(int) (*model.User, error)
	SendResetCode(email string) (string, error)
	ResetPassword(email, emailCode, password string) (string, error)
	DeleteUser(email string) error
}

type AddressRepository interface {
	CreateAddress(ad *model.Address) error
	GetAllAddresses() ([]*model.Address, error)
	DeleteAddress(id int) error
	UpdateAddress(p *model.Address) (string, error)
}
