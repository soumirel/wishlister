package service

type ServiceFactory interface {
	WishlistPermissionService() WishlistPermissionService
	UserService() UserService
	UserIdentityService() UserIdentityService
}
