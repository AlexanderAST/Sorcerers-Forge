package teststore

import (
	"Diploma/internal/model"
	_ "Diploma/internal/store"
	"Diploma/internal/store/sqlstore"
	sqlstore2 "Diploma/internal/store/sqlstore"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := sqlstore.New(db)

	mock.ExpectQuery("INSERT INTO users").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))

	user := &model.User{Email: "example@mail.ru", Password: "qwerty1245"}

	err = store.User().Create(user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if user.ID != 0 {
		t.Errorf("unexpected category ID: got %d, want %d", user.ID, 0)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore2.TestDB(t, "host=localhost user=postgres password=qwerty dbname=postgres port=5436 sslmode=disable")

	defer teardown("users")

	s := sqlstore2.New(db)
	email := "helloworeld@exame.org"

	u := model.TestUser(t)
	u.Email = email

	s.User().Create(u)

	u, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
