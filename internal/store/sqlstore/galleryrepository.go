package sqlstore

import (
	"Diploma/internal/model"
	"database/sql"
	"errors"
	"fmt"
)

type GalleryRepository struct {
	store *Store
}

func (r *GalleryRepository) CreateGallery(g *model.Gallery) error {
	return r.store.db.QueryRow("INSERT INTO gallery(photo, description) values ($1,$2) RETURNING id", g.Photo, g.Description).Scan(&g.ID)
}

func (r *GalleryRepository) DeleteGallery(id int) error {

	p := &model.Gallery{}

	if err := r.store.db.QueryRow("SELECT * FROM gallery WHERE id = $1", id).Scan(&p.ID, &p.Photo, &p.Description); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	if err := r.store.db.QueryRow("DELETE FROM gallery WHERE id = $1", id); err != nil {
		return err.Err()
	}
	return nil
}

func (r *GalleryRepository) GetAllGallery() ([]*model.Gallery, error) {

	rows, err := r.store.db.Query("SELECT * FROM gallery")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Gallery, 0)
	for rows.Next() {
		pd := new(model.Gallery)
		err := rows.Scan(&pd.ID, &pd.Photo, &pd.Description)
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

func (r *GalleryRepository) UpdateGallery(g *model.Gallery) (string, error) {

	if err := r.store.db.QueryRow("UPDATE gallery SET photo = $2, description= $3 WHERE id = $1", g.ID, g.Photo, g.Description).Scan(&g.ID, &g.Photo, &g.Description); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("No rows were returned!")
			return "ok", nil
		case err == nil:
			fmt.Println(g.ID)
		default:
			return "ok", err
		}

	}

	return "success", nil
}
