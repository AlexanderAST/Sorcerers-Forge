package handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type ctxKey int8

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(corsMiddleware)
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	//user
	s.router.HandleFunc("/sign-up", s.handleUsersCreate()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/sign-in", s.handleSessionsCreate()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/confirmCode", s.handleUsersConfirm()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteUsers", s.deleteUser()).Methods("OPTIONS", "DELETE")
	s.router.HandleFunc("/resetCode", s.sendResetCode()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/resetPassword", s.resetPassword()).Methods("OPTIONS", "POST")

	//jwt
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.userIdentity)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("OPTIONS", "GET")
	//catalog
	s.router.HandleFunc("/createCategory", s.createCategory()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteCategory", s.deleteCategory()).Methods("OPTIONS", "DELETE")
	s.router.HandleFunc("/categories", s.getAllCategory()).Methods("OPTIONS", "GET")
	s.router.HandleFunc("/createProduct", s.createProduct()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteProduct", s.deleteProduct()).Methods("OPTIONS", "DELETE")
	s.router.HandleFunc("/productsCategory", s.getAllProductsCategories()).Methods("OPTIONS", "GET")
	s.router.HandleFunc("/findProductByID", s.findProductById()).Methods("OPTIONS", "GET")
	products := s.router.PathPrefix("/products").Subrouter()
	products.Use(s.optionalUserIdentity)
	products.HandleFunc("", s.getAllProductsWrapper).Methods("OPTIONS", "GET")

	cart := s.router.PathPrefix("/cart").Subrouter()
	cart.Use(s.userIdentity)
	cart.HandleFunc("/createCart", s.createCartItem()).Methods("OPTIONS", "POST")
	cart.HandleFunc("/getUserCart", s.getAllCartUsers()).Methods("OPTIONS", "GET")
	cart.HandleFunc("/deleteCart", s.deleteCartItem()).Methods("OPTIONS", "DELETE")
	cart.HandleFunc("/createOrder", s.sendUserOrder()).Methods("OPTIONS", "GET")
	favorite := s.router.PathPrefix("/favorite").Subrouter()
	favorite.Use(s.userIdentity)
	favorite.HandleFunc("/favorite", s.getAllFavoriteUsers()).Methods("OPTIONS", "GET")
	favorite.HandleFunc("/deleteFavorite", s.deleteFavoriteItem()).Methods("OPTIONS", "DELETE")
	favorite.HandleFunc("/likeProduct", s.handleLikes()).Methods("OPTIONS", "POST")
	favorite.HandleFunc("/addToFavorite", s.createFavoriteItem()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/updateProduct", s.updateProduct()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/getLikeCount", s.getAllLikes()).Methods("OPTIONS", "GET")

	//profile
	profile := s.router.PathPrefix("/profile").Subrouter()
	profile.Use(s.userIdentity)
	profile.HandleFunc("/createProfile", s.createProfile()).Methods("OPTIONS", "POST")
	profile.HandleFunc("/updateProfile", s.updateProfile()).Methods("OPTIONS", "POST")
	profile.HandleFunc("/deleteProfile", s.deleteProfile()).Methods("OPTIONS", "DELETE")
	profile.HandleFunc("/profile", s.takeProfile()).Methods("OPTIONS", "GET")
	profile.HandleFunc("/orderHistory", s.createOrderHistory()).Methods("OPTIONS", "GET")
	//send msg from help
	s.router.HandleFunc("/sendMsg", s.sendHelp()).Methods("OPTIONS", "POST")
	//reviews
	reviews := s.router.PathPrefix("").Subrouter()
	reviews.Use(s.userIdentity)
	reviews.HandleFunc("/createReview", s.createReview()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/getReviews", s.getAllReviews()).Methods("OPTIONS", "GET")
	reviews.HandleFunc("/deleteReview", s.deleteReview()).Methods("OPTIONS", "DELETE")
	reviews.HandleFunc("/updateReviews", s.updateReview()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/getReview", s.getAllReviewsByProduct()).Methods("OPTIONS", "GET")
	//gallery
	s.router.HandleFunc("/createGallery", s.createGallery()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteGallery", s.deleteGallery()).Methods("OPTIONS", "DELETE")
	s.router.HandleFunc("/getGallery", s.getAllGallery()).Methods("OPTIONS", "GET")
	s.router.HandleFunc("/updateGallery", s.updateGallery()).Methods("OPTIONS", "POST")
	//admin
	s.router.HandleFunc("/admin-sign-up", s.handleAdminsCreate()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/admin-sign-in", s.handleSessionsAdminsCreate()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteAdminUsers", s.deleteAdmin()).Methods("OPTIONS", "DELETE")
	s.router.HandleFunc("/resetAdminCode", s.sendResetCodeAdmins()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/resetAdminPassword", s.resetAdminsPassword()).Methods("OPTIONS", "POST")
	privateAdmin := s.router.PathPrefix("/privateAdmin").Subrouter()
	privateAdmin.Use(s.authenticateAdmins)
	privateAdmin.HandleFunc("/whoamiAdmin", s.handleAdminsWhoami()).Methods("OPTIONS", "GET")
	//address
	s.router.HandleFunc("/createAddress", s.createAddress()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/getAllAddress", s.getAllAddresses()).Methods("OPTIONS", "GET")
	s.router.HandleFunc("/deleteAddress", s.deleteAddress()).Methods("OPTIONS", "DELETE")
	s.router.HandleFunc("/updateAddress", s.updateAddress()).Methods("OPTIONS", "POST")
	//content
	s.router.HandleFunc("/uploadCatalogPhoto", s.uploadCatalogPhotos()).Methods("OPTIONS", "PUT")
	s.router.HandleFunc("/uploadProfilePhoto", s.uploadProfilePhotos()).Methods("OPTIONS", "PUT")
	s.router.HandleFunc("/uploadReviewsPhoto", s.uploadReviewsPhotos()).Methods("OPTIONS", "PUT")
	s.router.HandleFunc("/uploadGalleryPhoto", s.uploadGalleryPhotos()).Methods("OPTIONS", "PUT")
	s.router.HandleFunc("/uploadApks", s.uploadApks()).Methods("OPTIONS", "PUT")
	s.router.HandleFunc("/getApk", s.getApks()).Methods("OPTIONS", "GET")
	s.router.HandleFunc("/downloadApk", s.downloadApk()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteCatalogPhoto", s.deleteCatalogPhoto()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteApk", s.deleteApk()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteProfilePhoto", s.deleteProfilePhoto()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteReviewsPhoto", s.deleteReviewsPhoto()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/deleteGalleryPhoto", s.deleteGalleryPhoto()).Methods("OPTIONS", "POST")
	s.router.HandleFunc("/get-Photo", s.handleGetPhoto()).Methods("OPTIONS", "GET")
	//LearnRequest
	learn := s.router.PathPrefix("/learn").Subrouter()
	learn.Use(s.optionalUserIdentity)
	learn.HandleFunc("", s.learnWrapper).Methods("OPTIONS", "POST")
	//logs info
	logrus.Info("GET /static/") // Статические файлы
	logrus.Info("POST /sign-up")
	logrus.Info("POST /sign-in")
	logrus.Info("POST /confirmCode")
	logrus.Info("DELETE /deleteUsers")
	logrus.Info("POST /resetCode")
	logrus.Info("POST /resetPassword")

	logrus.Info("GET /private/whoami") // Аутентификация JWT

	logrus.Info("POST /createCategory")
	logrus.Info("DELETE /deleteCategory")
	logrus.Info("GET /categories")
	logrus.Info("POST /createProduct")
	logrus.Info("DELETE /deleteProduct")
	logrus.Info("GET /products")
	logrus.Info("POST /productsCategory")
	logrus.Info("GET /products/productsCategory")
	logrus.Info("POST /cart/createCart")
	logrus.Info("GET /cart/getUserCart")
	logrus.Info("DELETE /cart/deleteCart")
	logrus.Info("POST /cart/createOrder")
	logrus.Info("POST /favorite/favorite")
	logrus.Info("GET /favorite/favorite")
	logrus.Info("DELETE /favorite/deleteFavorite")
	logrus.Info("POST /favorite/likeProduct")
	logrus.Info("POST /favorite/addToFavorite")
	logrus.Info("POST /updateProduct")
	logrus.Info("GET /getLikeCount")
	logrus.Info("POST /profile/createProfile")
	logrus.Info("POST /profile/updateProfile")
	logrus.Info("DELETE /profile/deleteProfile")
	logrus.Info("GET /profile/profile")
	logrus.Info("GET /profile/orderHistory")
	logrus.Info("POST /sendMsg")
	logrus.Info("POST /createReview")
	logrus.Info("GET /getReviews")
	logrus.Info("DELETE /deleteReview")
	logrus.Info("POST /updateReviews")
	logrus.Info("POST /createGallery")
	logrus.Info("DELETE /deleteGallery")
	logrus.Info("GET /getGallery")
	logrus.Info("POST /updateGallery")
	logrus.Info("POST /admin-sign-up")
	logrus.Info("POST /admin-sign-in")
	logrus.Info("DELETE /deleteAdminUsers")
	logrus.Info("POST /resetAdminCode")
	logrus.Info("POST /resetAdminPassword")
	logrus.Info("GET /privateAdmin/whoamiAdmin")
	logrus.Info("POST /createAddress")
	logrus.Info("GET /getAllAddress")
	logrus.Info("DELETE /deleteAddress")
	logrus.Info("POST /updateAddress")
	logrus.Info("PUT /uploadCatalogPhoto")
	logrus.Info("PUT /uploadProfilePhoto")
	logrus.Info("PUT /uploadReviewsPhoto")
	logrus.Info("PUT /uploadGalleryPhoto")
	logrus.Info("PUT /uploadApks")
	logrus.Info("GET /getApk")
	logrus.Info("POST /downloadApk")
	logrus.Info("POST /deleteCatalogPhoto")
	logrus.Info("POST /deleteApk")
	logrus.Info("POST /deleteProfilePhoto")
	logrus.Info("POST /deleteReviewsPhoto")
	logrus.Info("POST /deleteGalleryPhoto")
}