package sqlstore

import (
	"Diploma/internal/model"
	"Diploma/internal/store"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"regexp"
	"strconv"
	"strings"
)

type CatalogRepository struct {
	store *Store
}

func (r *CatalogRepository) CreateCategory(c *model.Category) error {
	return r.store.db.QueryRow("INSERT INTO product_category(name) values ($1) RETURNING id", c.Name).Scan(&c.ID)
}

func (r *CatalogRepository) AddLikeProduct(l *model.Like) error {
	return r.store.db.QueryRow("INSERT INTO likes(user_id, product_id) values ($1,$2) RETURNING id", l.UserId, l.ProductID).Scan(&l.ID)
}

func (r *CatalogRepository) CheckLikes(user_id, product_id int) (bool, error) {
	var exists bool
	if err := r.store.db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND product_id = $2)", user_id, product_id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *CatalogRepository) GetProductLikes() ([]*model.Like, error) {
	rows, err := r.store.db.Query("SELECT product_id, COUNT(*) AS like_count FROM likes GROUP BY product_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Like, 0)
	for rows.Next() {
		pd := new(model.Like)
		err := rows.Scan(&pd.UserId, &pd.ProductID)
		if err != nil {
			return nil, err
		}
		pdk = append(pdk, pd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pdk, err
}

func (r *CatalogRepository) GetAllCategories() ([]*model.Category, error) {

	rows, err := r.store.db.Query("SELECT * FROM product_category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Category, 0)
	for rows.Next() {
		pd := new(model.Category)
		err := rows.Scan(&pd.ID, &pd.Name)
		if err != nil {
			return nil, err
		}
		pdk = append(pdk, pd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pdk, err
}

func (r *CatalogRepository) DeleteCategory(id int) error {

	u := &model.Category{}

	if err := r.store.db.QueryRow("SELECT * FROM product_category WHERE id = $1", id).Scan(&u.ID, &u.Name); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	if err := r.store.db.QueryRow("DELETE FROM product_category WHERE id = $1", id); err != nil {
		return err.Err()
	}
	return nil
}

func (r *CatalogRepository) CreateProduct(c *model.Products) error {

	return r.store.db.QueryRow("INSERT INTO products(name,description, price,reviews_mid, reviews_count, quantity, work_time, photo, category_id, is_active) VALUES ($1, $2, $3, $4, $5,$6,$7,$8,$9,$10) RETURNING id", c.Name, c.Description, c.Price, c.ReviewsMid, c.ReviewsCount, c.Quantity, c.WorkTime, c.Photo, c.CategoryID, c.IsActive).Scan(&c.ID)

}

func (r *CatalogRepository) DeleteProduct(id int) error {

	c := &model.Products{}

	if err := r.store.db.QueryRow("SELECT * FROM products WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Description, &c.Price, &c.ReviewsMid, &c.ReviewsCount, &c.Quantity, &c.WorkTime, &c.Photo, &c.CategoryID, &c.IsActive); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	_ = r.store.db.QueryRow("DELETE FROM favorite_items WHERE product_id = $1", id)

	_ = r.store.db.QueryRow("DELETE FROM cart_items WHERE product_id= $1", id)

	_ = r.store.db.QueryRow("DELETE FROM likes WHERE product_id= $1", id)

	if err := r.store.db.QueryRow("DELETE FROM products WHERE id = $1", id); err != nil {
		return err.Err()
	}
	return nil
}

func (r *CatalogRepository) GetAllProducts() ([]*model.Products, error) {

	rows, err := r.store.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Products, 0)
	for rows.Next() {
		pd := new(model.Products)
		err := rows.Scan(&pd.ID, &pd.Name, &pd.Description, &pd.Price, &pd.ReviewsMid, &pd.ReviewsCount, &pd.Quantity, &pd.WorkTime, &pd.Photo, &pd.CategoryID, &pd.IsActive)
		if err != nil {
			return nil, err
		}
		pdk = append(pdk, pd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pdk, err

}

func (r *CatalogRepository) GetAllProductsSortedByLikesDesc(offset, limit int) ([]*model.Products, error) {
	query := fmt.Sprintf(`
        SELECT p.*, COALESCE(l.like_count, 0) AS like_count
        FROM products p
        LEFT JOIN (
            SELECT product_id, COUNT(*) AS like_count
            FROM likes
            GROUP BY product_id
        ) l ON p.id = l.product_id
        ORDER BY like_count DESC
        LIMIT %d OFFSET %d`, limit, offset)

	rows, err := r.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows, true)
}

func (r *CatalogRepository) GetFilteredProductsSortedByLikesDesc(offset, limit int, filters map[string]interface{}) ([]*model.Products, error) {
	// Формируем базовый запрос
	query := fmt.Sprintf(`
        SELECT p.*, COALESCE(l.like_count, 0) AS like_count
        FROM products p
        LEFT JOIN (
            SELECT product_id, COUNT(*) AS like_count
            FROM likes
            GROUP BY product_id
        ) l ON p.id = l.product_id
        WHERE 1=1`)

	// Добавляем условия фильтрации, если они есть
	args := make([]interface{}, 0)
	for key, value := range filters {
		query += fmt.Sprintf(" AND %s = ?", key)
		args = append(args, value)
	}

	// Добавляем сортировку
	query += " ORDER BY like_count DESC"

	// Добавляем ограничение на количество возвращаемых результатов
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Выполняем запрос
	rows, err := r.store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Сканируем результаты
	return scanProducts(rows, true)
}

func (r *CatalogRepository) GetAllProductsSortedByPriceAsc(offset, limit int) ([]*model.Products, error) {
	query := fmt.Sprintf("SELECT * FROM products ORDER BY price ASC LIMIT %d OFFSET %d", limit, offset)
	rows, err := r.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows, false)
}

func (r *CatalogRepository) GetAllProductsSortedByPriceDesc(offset, limit int) ([]*model.Products, error) {
	query := fmt.Sprintf("SELECT * FROM products ORDER BY price DESC LIMIT %d OFFSET %d", limit, offset)
	rows, err := r.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows, false)
}

func (r *CatalogRepository) GetAllIsActive(offset, limit int, active string) ([]*model.Products, error) {
	query := fmt.Sprintf("SELECT * FROM products WHERE is_active = %s LIMIT %d OFFSET %d", active, limit, offset)
	rows, err := r.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows, false)
}

func (r *CatalogRepository) GetAllProductsPaginated(offset, limit int) ([]*model.Products, error) {
	query := fmt.Sprintf("SELECT * FROM products LIMIT %d OFFSET %d", limit, offset)
	rows, err := r.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows, false)
}

func scanProducts(rows *sql.Rows, withLikes bool) ([]*model.Products, error) {
	pdk := make([]*model.Products, 0)
	for rows.Next() {
		pd := new(model.Products)
		var likeCount int
		var err error
		if withLikes {
			err = rows.Scan(&pd.ID, &pd.Name, &pd.Description, &pd.Price, &pd.ReviewsMid, &pd.ReviewsCount, &pd.Quantity, &pd.WorkTime, &pd.Photo, &pd.CategoryID, &pd.IsActive, &likeCount)
		} else {
			err = rows.Scan(&pd.ID, &pd.Name, &pd.Description, &pd.Price, &pd.ReviewsMid, &pd.ReviewsCount, &pd.Quantity, &pd.WorkTime, &pd.Photo, &pd.CategoryID, &pd.IsActive)
		}
		if err != nil {
			return nil, err
		}
		pdk = append(pdk, pd)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pdk, nil
}

func (r *CatalogRepository) GetAllProductsSortedByLikesDesc1(offset, limit int, priceFilter, activeFilter, typeFilter string) ([]*model.Products, error) {
	query := fmt.Sprintf(`
        SELECT p.*, COALESCE(l.like_count, 0) AS like_count
        FROM products p
        LEFT JOIN (
            SELECT product_id, COUNT(*) AS like_count
            FROM likes
            GROUP BY product_id
        ) l ON p.id = l.product_id
        WHERE 1=1
		%s
        ORDER BY like_count DESC
        LIMIT %d OFFSET %d`, r.buildFilterConditions(priceFilter, activeFilter, typeFilter), limit, offset)

	rows, err := r.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows, true)
}

func (r *CatalogRepository) buildFilterConditions(priceFilter, activeFilter, typeFilter string) string {
	conditions := ""

	if priceFilter != "" {
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

		conditions += fmt.Sprintf(" AND price BETWEEN %v AND %v", start, finish)
	}

	if activeFilter != "" {
		conditions += fmt.Sprintf(" AND is_active = '%s'", activeFilter)
	}

	if typeFilter != "" {
		var id string
		_ = r.store.db.QueryRow("SELECT id FROM product_category WHERE name = $1", typeFilter).Scan(&id)

		conditions += fmt.Sprintf(" AND category_id  = %s", id)
	}

	return conditions
}

func (r *CatalogRepository) GetAllProductsCategories(categoryID int) ([]*model.Products, error) {

	rows, err := r.store.db.Query("SELECT * FROM products WHERE category_id = $1", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Products, 0)
	for rows.Next() {
		pd := new(model.Products)
		err := rows.Scan(&pd.ID, &pd.Name, &pd.Description, &pd.Price, &pd.ReviewsMid, &pd.ReviewsCount, &pd.Quantity, &pd.WorkTime, &pd.Photo, &pd.CategoryID, &pd.IsActive)
		if err != nil {
			return nil, err
		}
		pdk = append(pdk, pd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pdk, err
}

func (r *CatalogRepository) FindByID(id int) (*model.Products, error) {

	c := &model.Products{}

	if err := r.store.db.QueryRow("SELECT * FROM products WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Description, &c.Price, &c.ReviewsMid, &c.ReviewsCount, &c.Quantity, &c.WorkTime, &c.Photo, &c.CategoryID, &c.IsActive); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return c, nil
}

func (r *CatalogRepository) AddToCartProduct(cart *model.CartItem) error {
	var result string
	_ = r.store.db.QueryRow("SELECT photo from products WHERE id = $1", cart.ProductID).Scan(&result)

	return r.store.db.QueryRow("INSERT INTO cart_items(user_id, product_id, count, photo) values ($1,$2,$3,$4) RETURNING id", cart.UserId, cart.ProductID, cart.Count, result).Scan(&cart.ID)
}

func (r *CatalogRepository) GetAllCartUser(id int) ([]*model.Cart, error) {

	rows, err := r.store.db.Query("SELECT * FROM cart_items WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Cart, 0)
	for rows.Next() {
		pd := new(model.Cart)
		err := rows.Scan(&pd.ID, &pd.UserId, &pd.ProductID, &pd.Count, &pd.Photo)
		if err != nil {
			return nil, err
		}

		productName, err := r.store.Catalog().FindByID(pd.ProductID)

		if err != nil {
			return nil, err
		}

		userEmail, err := r.store.User().Find(pd.UserId)
		if err != nil {
			return nil, err
		}

		pd.UserEmail = userEmail.Email
		pd.ProductName = productName.Name
		pd.Price = productName.Price * pd.Count
		pd.Photo = productName.Photo
		pdk = append(pdk, pd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pdk, err
}

func (r *CatalogRepository) DeleteCart(productid, userID int) error {

	c := &model.CartItem{}

	if err := r.store.db.QueryRow("SELECT * FROM cart_items WHERE product_id = $1 AND user_id = $2", productid, userID).Scan(&c.ID, &c.UserId, &c.ProductID, &c.Count, &c.Photo); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	if err := r.store.db.QueryRow("DELETE FROM cart_items WHERE product_id = $1 AND user_id = $2", productid, userID); err != nil {
		return err.Err()
	}
	return nil
}

func (r *CatalogRepository) DeleteAllFromCart(userID int) error {
	if err := r.store.db.QueryRow("DELETE FROM cart_items WHERE user_id = $1", userID); err != nil {
		return err.Err()
	}
	return nil
}

func (r *CatalogRepository) AddToFavoriteProduct(cart *model.FavoriteItem) error {
	var result string
	_ = r.store.db.QueryRow("SELECT photo from products WHERE id = $1", cart.ProductID).Scan(&result)

	return r.store.db.QueryRow("INSERT INTO favorite_items(user_id, product_id, photo) values ($1,$2,$3) RETURNING id", cart.UserId, cart.ProductID, result).Scan(&cart.ID)
}

func (r *CatalogRepository) GetAllFavoriteUser(id int) ([]*model.Favorite, error) {

	rows, err := r.store.db.Query("SELECT * FROM favorite_items WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Favorite, 0)
	for rows.Next() {
		pd := new(model.Favorite)
		err := rows.Scan(&pd.ID, &pd.UserId, &pd.ProductID, &pd.Photo)
		if err != nil {
			return nil, err
		}

		productName, err := r.store.Catalog().FindByID(pd.ProductID)

		if err != nil {
			return nil, err
		}

		userEmail, err := r.store.User().Find(pd.UserId)
		if err != nil {
			return nil, err
		}

		pd.UserEmail = userEmail.Email
		pd.ProductName = productName.Name
		pd.Price = productName.Price
		pd.Photo = productName.Photo
		pdk = append(pdk, pd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pdk, err
}

func (r *CatalogRepository) DeleteFavorite(productid, userID int) error {

	c := &model.FavoriteItem{}

	if err := r.store.db.QueryRow("SELECT * FROM favorite_items WHERE product_id = $1 AND user_id = $2", productid, userID).Scan(&c.ID, &c.UserId, &c.ProductID, &c.Photo); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	if err := r.store.db.QueryRow("DELETE FROM favorite_items WHERE product_id = $1 AND user_id = $2", productid, userID); err != nil {

		if err := r.store.db.QueryRow("DELETE FROM likes WHERE product_id = $1 AND user_id = $2", productid, userID); err != nil {
			return err.Err()
		}
		return err.Err()
	}

	return nil
}
func (r *CatalogRepository) UpdateProduct(c *model.Products) (string, error) {

	query := "UPDATE products SET "
	args := []interface{}{c.ID}
	argCounter := 2

	if c.Name != "" {
		query += "name = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, c.Name)
		argCounter++
	}
	if c.Description != "" {
		query += "description = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, c.Description)
		argCounter++
	}
	if c.Price != 0 {
		query += "price = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, c.Price)
		argCounter++
	}
	if c.Quantity != 0 {
		query += "quantity = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, c.Quantity)
		argCounter++
	}

	if c.WorkTime != 0 {
		query += "work_time = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, c.WorkTime)
		argCounter++
	}
	if c.Photo != "" {
		query += "photo = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, c.Photo)
		argCounter++
	}
	if c.CategoryID != 0 {
		query += "category_id = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, c.CategoryID)
		argCounter++
	}
	if c.IsActive {
		query += "is_active = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, c.IsActive)
		argCounter++
	}

	query = strings.TrimSuffix(query, ", ")

	query += " WHERE id = $1 RETURNING id, name, description, price, quantity, work_time, photo, category_id, is_active"

	err := r.store.db.QueryRow(query, args...).Scan(&c.ID, &c.Name, &c.Description, &c.Price, &c.Quantity, &c.WorkTime, &c.Photo, &c.CategoryID, &c.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("No rows were returned!")
			return "ok", nil
		case err == nil:
			fmt.Println(c.ID)
		default:
			return "ok", err
		}
	}

	return "success", nil
}

func (r *CatalogRepository) GetProductsBetweenPrice(startPrice, secondPrice int, offset, limit int) ([]*model.Products, error) {
	query := fmt.Sprintf("SELECT * FROM products where price between %v and %v LIMIT %d OFFSET %d ", startPrice, secondPrice, limit, offset)
	rows, err := r.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows, false)
}

func (r *CatalogRepository) FilterByType(offset, limit int, category string) ([]*model.Products, error) {
	var id string
	err := r.store.db.QueryRow("SELECT id FROM product_category WHERE name = $1 LIMIT 1 OFFSET $2", category, offset).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Выполнение запроса для получения продуктов по идентификатору категории
	query := fmt.Sprintf("SELECT * FROM products WHERE category_id = '%s' OFFSET $1 LIMIT $2", id)
	rows, err := r.store.db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows, false)

}

func (r *CatalogRepository) CountFilteredProducts(priceFilter, activeFilter, typeFilter string) (int, error) {
	query := "SELECT COUNT(*) FROM products WHERE 1=1" + r.buildFilterConditions(priceFilter, activeFilter, typeFilter)

	var totalCount int
	err := r.store.db.QueryRow(query).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (r *CatalogRepository) GetUserCartProducts(userID int) ([]*model.Products, error) {
	query := `
        SELECT p.*
        FROM products p
        JOIN cart_items ci ON p.id = ci.product_id
        WHERE ci.user_id = $1
    `
	rows, err := r.store.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows, false)
}

func (r *CatalogRepository) GetUserFavoriteProducts(userID int) ([]*model.Products, error) {
	query := `
        SELECT p.*
        FROM products p
        JOIN favorite_items fi ON p.id = fi.product_id
        WHERE fi.user_id = $1
    `
	rows, err := r.store.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows, false)
}

func (r *CatalogRepository) IsProductInCart(userID, productID int) (bool, error) {
	var exists bool
	err := r.store.db.QueryRow("SELECT EXISTS(SELECT 1 FROM cart_items WHERE user_id = $1 AND product_id = $2)", userID, productID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *CatalogRepository) IsProductInFavorites(userID, productID int) (bool, error) {
	var exists bool
	err := r.store.db.QueryRow("SELECT EXISTS(SELECT 1 FROM favorite_items WHERE user_id = $1 AND product_id = $2)", userID, productID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *CatalogRepository) GetAllProductsWithFlags(userID int) ([]*model.ProductWithFlags, error) {
	products, err := r.GetAllProducts()
	if err != nil {
		return nil, err
	}

	// Получаем список товаров в корзине пользователя
	cartItems, err := r.GetAllCartUser(userID)
	if err != nil {
		return nil, err
	}
	cartProductIDs := make(map[int]bool)
	for _, item := range cartItems {
		cartProductIDs[item.ProductID] = true
	}

	// Получаем список избранных товаров пользователя
	favoriteItems, err := r.GetAllFavoriteUser(userID)
	if err != nil {
		return nil, err
	}
	favoriteProductIDs := make(map[int]bool)
	for _, item := range favoriteItems {
		favoriteProductIDs[item.ProductID] = true
	}

	// Добавляем флаги is_favorite и is_cart к каждому продукту
	productsWithFlags := make([]*model.ProductWithFlags, len(products))
	for i, product := range products {
		isFavorite := favoriteProductIDs[product.ID]
		isCart := cartProductIDs[product.ID]
		productsWithFlags[i] = &model.ProductWithFlags{
			Products:   product,
			IsFavorite: isFavorite,
			IsCart:     isCart,
		}
	}

	return productsWithFlags, nil
}

func (r *CatalogRepository) CreateOrder(cartItems []*model.Cart) (int, error) {
	if len(cartItems) == 0 {
		return 0, errors.New("empty cart items")
	}
	var totalSum int
	var productIDs []int
	var productCount []int

	for _, cart := range cartItems {
		totalSum += cart.Price
		productIDs = append(productIDs, cart.ProductID)
		productCount = append(productCount, cart.Count)
	}

	query := "INSERT INTO orders (user_id, product_id,product_count, summ) VALUES ($1,$2,$3,$4) RETURNING id"
	var orderID int
	if err := r.store.db.QueryRow(query, cartItems[0].UserId, pq.Array(productIDs), pq.Array(productCount), totalSum).Scan(&orderID); err != nil {
		return 0, err
	}

	return orderID, nil
}

func (r *CatalogRepository) GetUserOrderHistory(userID int) ([]*model.Orders, error) {
	var orders []*model.Orders

	query := "SELECT id,user_id, product_id,product_count, summ FROM orders WHERE user_id = $1"
	rows, err := r.store.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := &model.Orders{}
		var productIDs pq.Int64Array
		var productCounts pq.Int64Array

		if err := rows.Scan(&order.ID, &order.UserID, &productIDs, &productCounts, &order.Summ); err != nil {
			return nil, err
		}
		order.ProductID = make([]int, len(productIDs))
		order.ProductCount = make([]int, len(productCounts))
		products := make([]model.ProductCartInfo, len(order.ProductID))

		for i, v := range productIDs {
			order.ProductID[i] = int(v)
		}
		for i, v := range productCounts {
			order.ProductCount[i] = int(v)
		}

		for i := range products {
			products[i] = model.ProductCartInfo{
				ProductId: order.ProductID[i],
				Count:     order.ProductCount[i],
			}
		}

		order.Products = products

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *CatalogRepository) GetUserOrderHistoryByReviews(userID, productID int) (bool, error) {
	query := "SELECT user_id, product_id FROM orders WHERE user_id = $1 AND  $2 = ANY(product_id)"
	rows, err := r.store.db.Query(query, userID, productID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Если есть хотя бы одна строка в результатах запроса, значит, у пользователя есть заказ для этого продукта
	return rows.Next(), nil
}
