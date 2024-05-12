package sqlstore

import (
	"Diploma/internal/model"
	"Diploma/internal/store"
	"database/sql"
	"errors"
	"fmt"
)

type ProfileRepository struct {
	store *Store
}

func (r *ProfileRepository) CreateProfile(p *model.Profile) error {
	return r.store.db.QueryRow("INSERT INTO profile(user_id, name, surname, patronymic, contact, photo) values ($1, $2, $3, $4, $5, $6) RETURNING id", p.UserID, p.Name, p.Surname, p.Patronymic, p.Contact, p.Photo).Scan(&p.ID)
}

func (r *ProfileRepository) UpdateProfile(p *model.Profile) (string, error) {

	if err := r.store.db.QueryRow("UPDATE profile SET name= $2, surname = $3, patronymic=$4, contact=$5, photo= $6 WHERE user_id = $1", p.UserID, p.Name, p.Surname, p.Patronymic, p.Contact, p.Photo).Scan(&p.UserID, &p.Name, &p.Surname, &p.Patronymic, &p.Contact, &p.Photo); err != nil {
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

func (r *ProfileRepository) DeleteProfile(id int) error {

	p := &model.Profile{}

	if err := r.store.db.QueryRow("SELECT * FROM profile WHERE user_id = $1", id).Scan(&p.ID, &p.UserID, &p.Name, &p.Surname, &p.Patronymic, &p.Contact, &p.Photo); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	if err := r.store.db.QueryRow("DELETE FROM profile WHERE user_id = $1", id); err != nil {
		return err.Err()
	}
	return nil
}

func (r *ProfileRepository) FindByID(id int) (*model.Profile, error) {

	p := &model.Profile{}

	if err := r.store.db.QueryRow("SELECT * FROM profile WHERE  user_id = $1", id).Scan(&p.ID, &p.UserID, &p.Name, &p.Surname, &p.Patronymic, &p.Contact, &p.Photo); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return p, nil
}
