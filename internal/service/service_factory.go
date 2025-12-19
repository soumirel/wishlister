package service

import (
	"wishlister/internal/domain/repository"
	"wishlister/internal/domain/service"
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
