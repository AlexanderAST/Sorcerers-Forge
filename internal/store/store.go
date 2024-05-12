package store

type Store interface {
	User() UserRepository
	Catalog() CatalogRepository
	Profile() ProfileRepository
	Reviews() ReviewsRepository
	Gallery() GalleryRepository
	Admin() AdminRepository
	Address() AddressRepository
}
