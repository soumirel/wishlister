package service

type ServiceFactory interface {
	WishlistPermissionService() WishlistPermissionService
}
