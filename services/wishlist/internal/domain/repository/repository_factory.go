package repository

type RepositoryFactory interface {
	UserRepository() UserRepository
	WishRepository() WishRepository
	WishlistRepository() WishlistRepository
	WishlistPermissionRepository() WishlistPermissionRepository
	UserIdentityRepository() UserIdentityRepository
}
