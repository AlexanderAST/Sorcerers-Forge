package sqlstore

import (
	"Diploma/internal/model"
	"database/sql"
	"errors"
	"fmt"
)

type ReviewsRepository struct {
	store *Store
}

func (r *ReviewsRepository) CreateReviews(p *model.Reviews) error {

	result := r.store.db.QueryRow("INSERT INTO reviews(product_id, user_id, stars, message) values ($1, $2, $3, $4) RETURNING id", p.ProductId, p.UserID, p.Stars, p.Message).Scan(&p.ID)

	_, err := r.store.db.Exec(`
        UPDATE products
        SET reviews_count = subquery.reviews_count,
            reviews_mid = subquery.reviews_mid
        FROM (
            SELECT product_id, COUNT(*) AS reviews_count, ROUND(AVG(stars), 2) AS reviews_mid
            FROM reviews
            GROUP BY product_id
        ) subquery
        WHERE products.id = subquery.product_id;
    `)

	if err != nil {
		return err
	}
	return result
}

func (r *ReviewsRepository) GetAllReviews() ([]*model.Reviews, error) {

	_, err := r.store.db.Exec(`
        UPDATE products
        SET reviews_count = subquery.reviews_count,
            reviews_mid = subquery.reviews_mid
        FROM (
            SELECT product_id, COUNT(*) AS reviews_count, ROUND(AVG(stars), 2) AS reviews_mid
            FROM reviews
            GROUP BY product_id
        ) subquery
        WHERE products.id = subquery.product_id;
    `)
	if err != nil {
		return nil, err
	}

	rows, err := r.store.db.Query("SELECT * FROM reviews")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Reviews, 0)
	for rows.Next() {
		pd := new(model.Reviews)
		err := rows.Scan(&pd.ID, &pd.ProductId, &pd.UserID, &pd.Stars, &pd.Message)
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

func (r *ReviewsRepository) UpdateReviews(p *model.Reviews) (string, error) {
	if err := r.store.db.QueryRow("UPDATE reviews SET stars= $3, message = $4 WHERE product_id = $1 AND user_id = $2", p.ProductId, p.UserID, p.Stars, p.Message).Scan(&p.ProductId, &p.UserID, &p.Stars, &p.Message); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("No rows were returned!")
			return "ok", nil
		case err == nil:
			fmt.Println(p.ID)
		default:
			return "ok", err
		}

	}

	return "success", nil
}

func (r *ReviewsRepository) DeleteReview(productID, userID int) error {

	p := &model.Reviews{}

	if err := r.store.db.QueryRow("SELECT * FROM reviews WHERE product_id = $1 AND user_id=$2", productID, userID).Scan(&p.ID, &p.ProductId, &p.UserID, &p.Stars, &p.Message); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	if err := r.store.db.QueryRow("DELETE FROM reviews WHERE product_id = $1 AND user_id = $2", productID, userID); err != nil {
		return err.Err()
	}
	return nil
}

func (r *ReviewsRepository) GetAllReviewsFromProduct(productId int) ([]*model.Reviews, error) {

	_, err := r.store.db.Exec(`
        UPDATE products
        SET reviews_count = subquery.reviews_count,
            reviews_mid = subquery.reviews_mid
        FROM (
            SELECT product_id, COUNT(*) AS reviews_count, ROUND(AVG(stars), 2) AS reviews_mid
            FROM reviews
            GROUP BY product_id
        ) subquery
        WHERE products.id = subquery.product_id;
    `)
	if err != nil {
		return nil, err
	}

	rows, err := r.store.db.Query("SELECT * FROM reviews WHERE product_id = $1", productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Reviews, 0)
	for rows.Next() {
		pd := new(model.Reviews)
		err := rows.Scan(&pd.ID, &pd.ProductId, &pd.UserID, &pd.Stars, &pd.Message)
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
