package sqlstore

import (
	"Diploma/internal/model"
	"database/sql"
	"errors"
	"fmt"
)

type AddressRepository struct {
	store *Store
}

func (r *AddressRepository) CreateAddress(ad *model.Address) error {
	return r.store.db.QueryRow("INSERT INTO addresses(name,latlng) VALUES($1,$2) RETURNING id", ad.Name, ad.Latlng).Scan(&ad.ID)
}

func (r *AddressRepository) GetAllAddresses() ([]*model.Address, error) {

	rows, err := r.store.db.Query("SELECT * FROM addresses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pdk := make([]*model.Address, 0)
	for rows.Next() {
		pd := new(model.Address)
		err := rows.Scan(&pd.ID, &pd.Name, &pd.Latlng)
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

func (r *AddressRepository) DeleteAddress(id int) error {

	p := &model.Address{}

	if err := r.store.db.QueryRow("SELECT * FROM addresses WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Latlng); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	if err := r.store.db.QueryRow("DELETE FROM addresses WHERE id = $1", id); err != nil {
		return err.Err()
	}
	return nil
}

func (r *AddressRepository) UpdateAddress(p *model.Address) (string, error) {
	if err := r.store.db.QueryRow("UPDATE addresses SET name=$2,latlng=$3  WHERE id = $1", p.ID, p.Name, p.Latlng).Scan(&p.ID, &p.Name, &p.Latlng); err != nil {
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
