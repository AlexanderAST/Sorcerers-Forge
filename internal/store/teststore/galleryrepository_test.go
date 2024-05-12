package teststore

import (
	"Diploma/internal/model"
	"Diploma/internal/store/sqlstore"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

func TestCreateGallery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := sqlstore.New(db)

	gallery := &model.Gallery{
		Photo:       "photo.jpg",
		Description: "Крутое фото катаны нашей",
	}

	mock.ExpectQuery("INSERT INTO gallery").
		WithArgs(gallery.Photo, gallery.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = store.Gallery().CreateGallery(gallery)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if gallery.ID != 1 {
		t.Errorf("unexpected gallery ID: got %d, want %d", gallery.ID, 1)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllGallery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := sqlstore.New(db)

	expectedGallery := []*model.Gallery{
		{ID: 1, Photo: "photo1.jpg", Description: "Description 1"},
		{ID: 2, Photo: "photo2.jpg", Description: "Description 2"},
	}

	mock.ExpectQuery("SELECT \\* FROM gallery").
		WillReturnRows(sqlmock.NewRows([]string{"id", "photo", "description"}).
			AddRow(expectedGallery[0].ID, expectedGallery[0].Photo, expectedGallery[0].Description).
			AddRow(expectedGallery[1].ID, expectedGallery[1].Photo, expectedGallery[1].Description))

	gallery, err := store.Gallery().GetAllGallery()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(gallery, expectedGallery) {
		t.Errorf("unexpected gallery: got %+v, want %+v", gallery, expectedGallery)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
