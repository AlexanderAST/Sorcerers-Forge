package teststore

import (
	"Diploma/internal/model"
	_ "Diploma/internal/store"
	"Diploma/internal/store/sqlstore"
	_ "Diploma/internal/store/sqlstore"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCreateProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := sqlstore.New(db)

	profile := &model.Profile{
		UserID:     1,
		Name:       "Александр",
		Surname:    "Петров",
		Patronymic: "Александрович",
		Contact:    "example@example.com",
		Photo:      "profile.jpg",
	}

	mock.ExpectQuery("INSERT INTO profile").
		WithArgs(profile.UserID, profile.Name, profile.Surname, profile.Patronymic, profile.Contact, profile.Photo).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = store.Profile().CreateProfile(profile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if profile.ID != 1 {
		t.Errorf("unexpected profile ID: got %d, want %d", profile.ID, 1)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProfileFindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := sqlstore.New(db)

	expectedProfile := &model.Profile{
		ID:         1,
		UserID:     1,
		Name:       "Александр",
		Surname:    "Петров",
		Patronymic: "Александрович",
		Contact:    "example@example.com",
		Photo:      "profile.jpg",
	}

	mock.ExpectQuery("SELECT \\* FROM profile WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "surname", "patronymic", "contact", "photo"}).
			AddRow(expectedProfile.ID, expectedProfile.UserID, expectedProfile.Name, expectedProfile.Surname, expectedProfile.Patronymic, expectedProfile.Contact, expectedProfile.Photo))

	profile, err := store.Profile().FindByID(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(profile, expectedProfile) {
		t.Errorf("unexpected profile: got %+v, want %+v", profile, expectedProfile)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
