package wishlist

import (
	"context"
	"wishlister/internal/domain/entity"
	"wishlister/internal/domain/repository"
	"wishlister/internal/service"
	"wishlister/internal/usecase"
)

type WishlistUsecase struct {
	uofFactory usecase.UnitOfWorkFactory
}

func NewWishlistUsecase(uofFactory usecase.UnitOfWorkFactory) *WishlistUsecase {
	return &WishlistUsecase{
		uofFactory: uofFactory,
	}
}

func (uc *WishlistUsecase) GetWishlists(ctx context.Context, cmd GetWishlistsCommand) ([]entity.Wishlist, error) {
	uof := uc.uofFactory.NewUnitOfWork(false)
	var wishlistsRes []entity.Wishlist
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()

		persmissions, err := permissionService.GetPermissionsToWishlists(ctx, cmd.RequestorUserID)
		if err != nil {
			return err
		}
		wishlistsIds := persmissions.GetWishlitsIdsForAction(entity.ReadWishlistAction)
		if len(wishlistsIds) == 0 {
			return nil
		}

		wishlistRepo := rf.WishlistRepository()
		wishlists, err := wishlistRepo.GetWishlists(ctx, wishlistsIds)
		if err != nil {
			return err
		}
		wishlistsRes = make([]entity.Wishlist, 0, len(wishlists))
		for _, wl := range wishlists {
			wishlistsRes = append(wishlistsRes, *wl)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return wishlistsRes, nil
}

func (uc *WishlistUsecase) GetWishlist(ctx context.Context, cmd GetWishlistCommand) (entity.Wishlist, error) {
	uof := uc.uofFactory.NewUnitOfWork(false)
	var wishlistRes entity.Wishlist
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err := permissionService.CheckReadWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}
		wishlistRepo := rf.WishlistRepository()
		wishlist, err := wishlistRepo.GetWishlist(ctx, cmd.WishlistID)
		if err != nil {
			return err
		}
		if wishlist == nil {
			return entity.ErrWishDoesNotExist
		}
		wishlistRes = *wishlist
		return nil
	})
	if err != nil {
		return entity.Wishlist{}, err
	}
	return wishlistRes, nil
}

func (uc *WishlistUsecase) CreateWishlist(ctx context.Context, cmd CreateWishlistCommand) (entity.Wishlist, error) {
	wishlist := entity.NewWishlist(cmd.RequestorUserID, cmd.Name)
	uof := uc.uofFactory.NewUnitOfWork(true)
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		wishlistRepo := rf.WishlistRepository()
		err := wishlistRepo.CreateWishlist(ctx, wishlist)
		if err != nil {
			return err
		}
		persmission := entity.NewWishlistPersmission(
			cmd.RequestorUserID, wishlist.ID, entity.OwnerWishlistPermissionLevel,
		)
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err = permissionService.SaveWishlistPermission(ctx, persmission)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entity.Wishlist{}, err
	}
	return *wishlist, nil
}

func (uc *WishlistUsecase) UpdateWishlist(ctx context.Context, cmd UpdateWishlistCommand) (entity.Wishlist, error) {
	uof := uc.uofFactory.NewUnitOfWork(true)
	var wishlistRes entity.Wishlist
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err := permissionService.CheckModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}
		// TODO: To DAO
		wishlist := &entity.Wishlist{
			ID:   cmd.WishlistID,
			Name: cmd.Name,
		}
		wishlistRepo := rf.WishlistRepository()
		err = wishlistRepo.UpdateWishlist(ctx, wishlist)
		if err != nil {
			return err
		}
		wishlistRes = *wishlist
		return nil
	})
	if err != nil {
		return entity.Wishlist{}, err
	}
	return wishlistRes, nil
}

func (uc *WishlistUsecase) DeleteWishlist(ctx context.Context, cmd DeleteWishlistCommand) error {
	uof := uc.uofFactory.NewUnitOfWork(false)
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err := permissionService.CheckModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}
		wishlistRepo := rf.WishlistRepository()
		err = wishlistRepo.DeleteWishlist(ctx, cmd.WishlistID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
