package handler

import (
	"Diploma/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"sort"
	"strconv"
)

func (s *server) createCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		productCategory := &category{}

		if err := json.NewDecoder(r.Body).Decode(productCategory); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if productCategory.Name == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid name"))
			return
		}

		cat := &model.Category{
			Name: productCategory.Name,
		}

		if err := s.store.Catalog().CreateCategory(cat); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": cat.ID, "status": "success"})
	}
}

func (s *server) getAllCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := s.store.Catalog().GetAllCategories()
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"categories": categories})
	}
}

func (s *server) deleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		intId, _ := strconv.Atoi(id)

		if err := s.store.Catalog().DeleteCategory(intId); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "delete success"})
	}
}

func (s *server) createProduct() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		cat := &catalog{}

		if err := json.NewDecoder(r.Body).Decode(cat); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if cat.Name == "" || cat.Description == "" || cat.Price == 0 || cat.Photo == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid data"))
			return
		}

		c := &model.Products{
			Name:         cat.Name,
			Description:  cat.Description,
			Price:        cat.Price,
			ReviewsMid:   cat.ReviewsMid,
			ReviewsCount: cat.ReviewsCount,
			Quantity:     cat.Quantity,
			WorkTime:     cat.WorkTime,
			Photo:        cat.Photo,
			CategoryID:   cat.CategoryID,
			IsActive:     cat.IsActive,
		}

		if err := s.store.Catalog().CreateProduct(c); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": c.ID, "status": "success"})
	}
}

func (s *server) deleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		intId, _ := strconv.Atoi(id)

		if err := s.store.Catalog().DeleteProduct(intId); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "delete success"})
	}
}

func (s *server) getAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortParam := r.URL.Query().Get("sort")

		page := r.URL.Query().Get("page")
		pageSize := r.URL.Query().Get("pageSize")

		priceFilter := r.URL.Query().Get("price")

		activeFilter := r.URL.Query().Get("active")

		typeFilter := r.URL.Query().Get("type")

		pageInt, err := strconv.Atoi(page)

		userID, ok := r.Context().Value("userID").(int)

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err != nil || pageInt < 1 {
			pageInt = 1
		}

		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pageSizeInt < 1 {
			pageSizeInt = 10
		}

		var filteredProducts []*model.Products
		switch {
		case priceFilter != "":
			regex := *regexp.MustCompile("(\\d+)-(\\d+)")
			res := regex.FindAllStringSubmatch(priceFilter, -1)
			minPrice := ""
			maxPrice := ""
			for i := range res {
				minPrice += res[i][1]
				maxPrice += res[i][2]
			}

			start, _ := strconv.Atoi(minPrice)
			finish, _ := strconv.Atoi(maxPrice)

			filteredProducts, err = s.store.Catalog().GetProductsBetweenPrice(start, finish, (pageInt-1)*pageSizeInt, pageSizeInt)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		case activeFilter != "":
			filteredProducts, err = s.store.Catalog().GetAllIsActive((pageInt-1)*pageSizeInt, pageSizeInt, activeFilter)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		case typeFilter != "":
			filteredProducts, err = s.store.Catalog().FilterByType((pageInt-1)*pageSizeInt, pageSizeInt, typeFilter)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

		}
		if len(filteredProducts) > 0 {
			switch sortParam {
			case "price_asc":
				SortProductsByPriceAsc(filteredProducts)
			case "price_desc":
				SortProductsByPriceDesc(filteredProducts)
			case "popularity_desc":
				filteredProducts, err = s.store.Catalog().GetAllProductsSortedByLikesDesc1((pageInt-1)*pageSizeInt, pageSizeInt, priceFilter, activeFilter, typeFilter)
			}
		} else {
			switch sortParam {
			case "price_asc":
				filteredProducts, err = s.store.Catalog().GetAllProductsSortedByPriceAsc((pageInt-1)*pageSizeInt, pageSizeInt)
			case "price_desc":
				filteredProducts, err = s.store.Catalog().GetAllProductsSortedByPriceDesc((pageInt-1)*pageSizeInt, pageSizeInt)
			case "popularity_desc":
				filteredProducts, err = s.store.Catalog().GetAllProductsSortedByLikesDesc((pageInt-1)*pageSizeInt, pageSizeInt)
			default:
				filteredProducts, err = s.store.Catalog().GetAllProductsPaginated((pageInt-1)*pageSizeInt, pageSizeInt)
			}
		}

		var productsWithExtraInfo []*model.ProductWithFlags

		for _, product := range filteredProducts {
			isFavorite, err := s.store.Catalog().IsProductInFavorites(userID, product.ID)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			isInCart, err := s.store.Catalog().IsProductInCart(userID, product.ID)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			productWithExtraInfo := &model.ProductWithFlags{
				Products:   product,
				IsFavorite: isFavorite,
				IsCart:     isInCart,
			}
			productsWithExtraInfo = append(productsWithExtraInfo, productWithExtraInfo)
		}

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		totalCount, err := s.store.Catalog().CountFilteredProducts(priceFilter, activeFilter, typeFilter)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalPages := (totalCount + pageSizeInt - 1) / pageSizeInt

		s.respond(w, r, http.StatusOK, map[string]interface{}{
			"countPages":    totalPages,
			"countProducts": totalCount,
			"products":      productsWithExtraInfo,
			"page":          pageInt,
			"pageSize":      pageSizeInt,
		})
	}
}

func (s *server) getAllProductsNoAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortParam := r.URL.Query().Get("sort")

		page := r.URL.Query().Get("page")
		pageSize := r.URL.Query().Get("pageSize")

		priceFilter := r.URL.Query().Get("price")

		activeFilter := r.URL.Query().Get("active")

		typeFilter := r.URL.Query().Get("type")

		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			pageInt = 1
		}

		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pageSizeInt < 1 {
			pageSizeInt = 10
		}

		var filteredProducts []*model.Products
		switch {
		case priceFilter != "":
			regex := *regexp.MustCompile("(\\d+)-(\\d+)")
			res := regex.FindAllStringSubmatch(priceFilter, -1)
			minPrice := ""
			maxPrice := ""
			for i := range res {
				minPrice += res[i][1]
				maxPrice += res[i][2]
			}

			start, _ := strconv.Atoi(minPrice)
			finish, _ := strconv.Atoi(maxPrice)

			filteredProducts, err = s.store.Catalog().GetProductsBetweenPrice(start, finish, (pageInt-1)*pageSizeInt, pageSizeInt)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		case activeFilter != "":
			filteredProducts, err = s.store.Catalog().GetAllIsActive((pageInt-1)*pageSizeInt, pageSizeInt, activeFilter)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		case typeFilter != "":
			filteredProducts, err = s.store.Catalog().FilterByType((pageInt-1)*pageSizeInt, pageSizeInt, typeFilter)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

		}

		if len(filteredProducts) > 0 {
			switch sortParam {
			case "price_asc":
				SortProductsByPriceAsc(filteredProducts)
			case "price_desc":
				SortProductsByPriceDesc(filteredProducts)
			case "popularity_desc":
				filteredProducts, err = s.store.Catalog().GetAllProductsSortedByLikesDesc1((pageInt-1)*pageSizeInt, pageSizeInt, priceFilter, activeFilter, typeFilter)
			}
		} else {
			switch sortParam {
			case "price_asc":
				filteredProducts, err = s.store.Catalog().GetAllProductsSortedByPriceAsc((pageInt-1)*pageSizeInt, pageSizeInt)
			case "price_desc":
				filteredProducts, err = s.store.Catalog().GetAllProductsSortedByPriceDesc((pageInt-1)*pageSizeInt, pageSizeInt)
			case "popularity_desc":
				filteredProducts, err = s.store.Catalog().GetAllProductsSortedByLikesDesc((pageInt-1)*pageSizeInt, pageSizeInt)
			default:
				filteredProducts, err = s.store.Catalog().GetAllProductsPaginated((pageInt-1)*pageSizeInt, pageSizeInt)
			}
		}

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.Catalog().CountFilteredProducts(priceFilter, activeFilter, typeFilter)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalPages := (totalCount + pageSizeInt - 1) / pageSizeInt

		s.respond(w, r, http.StatusOK, map[string]interface{}{
			"countPages":    totalPages,
			"countProducts": totalCount,
			"products":      filteredProducts,
			"page":          pageInt,
			"pageSize":      pageSizeInt,
		})
	}
}

func SortProductsByPriceAsc(products []*model.Products) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].Price < products[j].Price
	})
}

func SortProductsByPriceDesc(products []*model.Products) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].Price > products[j].Price
	})
}

func (s *server) getAllProductsCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		intId, _ := strconv.Atoi(id)

		products, err := s.store.Catalog().GetAllProductsCategories(intId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"products": products})
	}
}

func (s *server) createCartItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cart := &cartItems{}

		if err := json.NewDecoder(r.Body).Decode(cart); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID := r.Context().Value("userID").(int)

		c := &model.CartItem{
			UserId:    userID,
			ProductID: cart.ProductID,
			Count:     cart.Count,
		}
		if err := s.store.Catalog().AddToCartProduct(c); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": c.ID, "status": "success"})
	}
}

func (s *server) getAllCartUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		carts, err := s.store.Catalog().GetAllCartUser(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"cart": carts})
	}
}

func (s *server) deleteCartItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		productID := r.URL.Query().Get("productId")
		idResult, _ := strconv.Atoi(productID)
		userID := r.Context().Value("userID").(int)

		if err := s.store.Catalog().DeleteCart(idResult, userID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "delete success"})
	}
}

func (s *server) findProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		idResult, _ := strconv.Atoi(id)

		name, err := s.store.Catalog().FindByID(idResult)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"name": name})
	}
}

func (s *server) createFavoriteItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cart := &favoriteItems{}

		if err := json.NewDecoder(r.Body).Decode(cart); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID := r.Context().Value("userID").(int)

		c := &model.FavoriteItem{
			UserId:    userID,
			ProductID: cart.ProductID,
		}

		exists, err := s.store.Catalog().CheckLikes(c.UserId, c.ProductID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if exists {
			s.error(w, r, http.StatusBadRequest, errors.New("user already liked this product"))
			return
		}

		if err := s.store.Catalog().AddToFavoriteProduct(c); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		addlike := &model.Like{

			UserId:    c.UserId,
			ProductID: c.ProductID,
		}

		if err := s.store.Catalog().AddLikeProduct(addlike); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": c.ID, "status": "success"})
	}
}

func (s *server) getAllFavoriteUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		favorite, err := s.store.Catalog().GetAllFavoriteUser(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"favorite": favorite})
	}
}

func (s *server) deleteFavoriteItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID := r.URL.Query().Get("productId")
		idResult, _ := strconv.Atoi(productID)
		userID := r.Context().Value("userID").(int)

		if err := s.store.Catalog().DeleteFavorite(idResult, userID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "delete success"})
	}
}

func (s *server) updateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &updateCatalog{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		existingProduct, err := s.store.Catalog().FindByID(input.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if input.Name != "" {
			existingProduct.Name = input.Name
		}
		if input.Description != "" {
			existingProduct.Description = input.Description
		}
		if input.Price != 0 {
			existingProduct.Price = input.Price
		}
		if input.Quantity != 0 {
			existingProduct.Quantity = input.Quantity
		}
		if input.WorkTime != 0 {
			existingProduct.WorkTime = input.WorkTime
		}
		if input.Photo != "" {
			existingProduct.Photo = input.Photo
		}
		if input.CategoryID != 0 {
			existingProduct.CategoryID = input.CategoryID
		}
		if input.IsActive {
			existingProduct.IsActive = input.IsActive
		}

		if _, err := s.store.Catalog().UpdateProduct(existingProduct); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": existingProduct.ID, "status": "success"})
	}
}

func (s *server) handleLikes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &reqLikes{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID := r.Context().Value("userID").(int)

		exists, err := s.store.Catalog().CheckLikes(userID, input.ProductID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if exists {
			s.error(w, r, http.StatusBadRequest, errors.New("user already liked this product"))
			return
		}

		addlike := &model.Like{

			UserId:    userID,
			ProductID: input.ProductID,
		}

		if err := s.store.Catalog().AddLikeProduct(addlike); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusCreated, map[string]interface{}{"status": "success", "id": addlike.ID})
	}
}

func (s *server) getAllLikes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		like, err := s.store.Catalog().GetProductLikes()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"gallery": like})
	}
}

func (s *server) getAllProductsWrapper(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if header == "" {
		s.getAllProductsNoAuth()(w, r)
		return
	}

	s.getAllProducts()(w, r)
}
