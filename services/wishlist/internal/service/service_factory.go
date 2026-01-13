package service

import (
	"github.com/soumirel/wishlister/wishlist/internal/domain/repository"
	"github.com/soumirel/wishlister/wishlist/internal/domain/service"
)

type serviceFactory struct {
	rf repository.RepositoryFactory
}

func NewServiceFactory(rf repository.RepositoryFactory) service.ServiceFactory {
	return &serviceFactory{
		rf: rf,
	}
}

func (f *serviceFactory) WishlistPermissionService() service.WishlistPermissionService {
	wishlistPermissionRepo := f.rf.WishlistPermissionRepository()
	return newWishlistPersmissionService(wishlistPermissionRepo)
}

func (f *serviceFactory) UserService() service.UserService {
	userRepo := f.rf.UserRepository()
	return newUserService(userRepo)
}

func (f *serviceFactory) UserIdentityService() service.UserIdentityService {
	userIdentityRepo := f.rf.UserIdentityRepository()
	return newUserIdentityServivce(userIdentityRepo)
}
