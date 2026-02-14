package service

type ServiceFactory interface {
	GetWishlistCoreReadService() WishlistCoreReadService
}
