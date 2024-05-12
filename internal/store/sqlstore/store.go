package sqlstore

import (
	"Diploma/internal/store"
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	db                *sql.DB
	userRepository    *UserRepository
	catalogRepository *CatalogRepository
	profileRepository *ProfileRepository
	reviewsRepository *ReviewsRepository
	galleryRepository *GalleryRepository
	adminRepository   *AdminRepository
	addressRepository *AddressRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {

	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{store: s}

	return s.userRepository
}

func (s *Store) Catalog() store.CatalogRepository {
	if s.catalogRepository != nil {
		return s.catalogRepository
	}

	s.catalogRepository = &CatalogRepository{store: s}

	return s.catalogRepository
}

func (s *Store) Profile() store.ProfileRepository {
	if s.profileRepository != nil {
		return s.profileRepository
	}

	s.profileRepository = &ProfileRepository{store: s}

	return s.profileRepository
}

func (s *Store) Reviews() store.ReviewsRepository {
	if s.reviewsRepository != nil {
		return s.reviewsRepository
	}

	s.reviewsRepository = &ReviewsRepository{store: s}

	return s.reviewsRepository
}

func (s *Store) Gallery() store.GalleryRepository {
	if s.galleryRepository != nil {
		return s.galleryRepository
	}

	s.galleryRepository = &GalleryRepository{store: s}

	return s.galleryRepository
}

func (s *Store) Admin() store.AdminRepository {

	if s.adminRepository != nil {
		return s.adminRepository
	}

	s.adminRepository = &AdminRepository{store: s}

	return s.adminRepository
}

func (s *Store) Address() store.AddressRepository {

	if s.addressRepository != nil {
		return s.addressRepository
	}

	s.addressRepository = &AddressRepository{store: s}

	return s.addressRepository
}
