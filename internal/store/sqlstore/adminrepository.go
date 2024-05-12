package sqlstore

import (
	"Diploma/configs"
	"Diploma/internal/model"
	"Diploma/internal/store"
	"database/sql"
	"errors"
)

type AdminRepository struct {
	store *Store
}

func (r *AdminRepository) Create(u *model.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow("INSERT INTO admins(email,encrypted_password) VALUES($1, $2) RETURNING id", u.Email, u.EncryptedPassword).Scan(&u.ID)

}

func (r *AdminRepository) FindByEmail(email string) (*model.User, error) {

	u := &model.User{}

	if err := r.store.db.QueryRow("SELECT id,email,encrypted_password FROM admins WHERE email = $1", email).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *AdminRepository) Find(id int) (*model.User, error) {

	u := &model.User{}

	if err := r.store.db.QueryRow("SELECT id,email,encrypted_password FROM admins WHERE id = $1", id).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *AdminRepository) SendResetCode(email string) (string, error) {

	u := &model.User{}

	if err := r.store.db.QueryRow("SELECT * FROM admins WHERE email = $1", email).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return "invalid email", err
		}
		return "", err
	}

	configs.RunSmtp(email)
	return "ok", nil
}

func (r *AdminRepository) ResetPassword(email, emailCode, password string) (string, error) {
	u := &model.User{}
	emCode := configs.GenToken(email)

	if emailCode == emCode {
		if err := r.store.db.QueryRow("UPDATE admins SET encrypted_password = $2 WHERE email = $1", email, password).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return "ok", nil
			default:
				return "ok", err
			}
		}
		return "ok", nil
	} else {
		return "bad", errors.New("invalid email code")
	}
}

func (r *AdminRepository) DeleteUser(email string) error {

	u := &model.User{}

	if err := r.store.db.QueryRow("SELECT * FROM admins WHERE email = $1", email).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return err
	}

	if err := r.store.db.QueryRow("DELETE FROM admins WHERE email = $1", email); err != nil {
		return err.Err()
	}
	return nil
}
