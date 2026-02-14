package service

import "github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"

type svcFactory struct {
	wishlistCoreReadSvc service.WishlistCoreReadService
}

func NewServiceFactory(
	wishlistCoreReadSvc service.WishlistCoreReadService,
) *svcFactory {
	return &svcFactory{
		wishlistCoreReadSvc: wishlistCoreReadSvc,
	}
}

func (f *svcFactory) GetWishlistCoreReadService() service.WishlistCoreReadService {
	return f.wishlistCoreReadSvc
}
