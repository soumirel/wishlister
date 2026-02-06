package wish

import (
	"context"
	"time"

	"github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/services/wishlist/internal/domain/repository"
	"github.com/soumirel/wishlister/services/wishlist/internal/service"
	"github.com/soumirel/wishlister/services/wishlist/internal/usecase"
)

type WishUsecase struct {
	uofFactory usecase.UnitOfWorkFactory
}

func NewWishUsecase(
	uofFactory usecase.UnitOfWorkFactory,
) *WishUsecase {
	return &WishUsecase{
		uofFactory: uofFactory,
	}
}

func (uc *WishUsecase) GetWish(ctx context.Context, cmd GetWishCommand) (entity.Wish, error) {
	uof := uc.uofFactory.NewUnitOfWork(false)
	var wishRes entity.Wish
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err := permissionService.CheckReadWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}
		wishRepo := rf.WishRepository()
		wish, err := wishRepo.GetWish(ctx, cmd.WishlistID, cmd.WishID)
		if err != nil {
			return err
		}
		if wish == nil {
			return entity.ErrWishlistDoesNotExist
		}
		wishRes = *wish
		return nil
	})
	if err != nil {
		return entity.Wish{}, err
	}
	return wishRes, nil
}

func (uc *WishUsecase) GetWishesFromWishlist(ctx context.Context, cmd GetWishesFromWishlistCommand) ([]entity.Wish, error) {
	uof := uc.uofFactory.NewUnitOfWork(true)
	var wishesRes []entity.Wish
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permSvc := svcFactory.WishlistPermissionService()
		err := permSvc.CheckReadWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}
		wishRepo := rf.WishRepository()
		wishes, err := wishRepo.GetWishesFromWishlist(ctx, cmd.WishlistID)
		if err != nil {
			return err
		}
		wishesRes = make([]entity.Wish, 0, len(wishes))
		for _, w := range wishes {
			wishesRes = append(wishesRes, *w)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return wishesRes, nil
}

func (uc *WishUsecase) CreateWish(ctx context.Context, cmd CreateWishCommand) (entity.Wish, error) {
	uof := uc.uofFactory.NewUnitOfWork(true)
	var wishRes entity.Wish
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err := permissionService.CheckModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}
		wishlistRepo := rf.WishlistRepository()
		wishlist, err := wishlistRepo.GetWishlist(ctx, cmd.WishlistID)
		if err != nil {
			return err
		}
		wish, err := wishlist.NewWish(cmd.WishName)
		if err != nil {
			return err
		}
		wishRepo := rf.WishRepository()
		err = wishRepo.CreateWish(ctx, wish)
		if err != nil {
			return err
		}
		wishRes = *wish
		return nil
	})
	if err != nil {
		return entity.Wish{}, err
	}
	return wishRes, nil
}

func (uc *WishUsecase) UpdateWish(ctx context.Context, cmd UpdateWishCommand) (entity.Wish, error) {
	uof := uc.uofFactory.NewUnitOfWork(true)
	var wishRes entity.Wish
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err := permissionService.CheckModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}
		wishRepo := rf.WishRepository()
		wish, err := wishRepo.GetWish(ctx, cmd.WishlistID, cmd.WishID)
		if err != nil {
			return err
		}
		wish.UpdateName(cmd.WishName)
		err = wishRepo.UpdateWish(ctx, wish)
		if err != nil {
			return err
		}
		wishRes = *wish
		return nil
	})
	if err != nil {
		return entity.Wish{}, err
	}
	return wishRes, nil
}

func (uc *WishUsecase) DeleteWish(ctx context.Context, cmd DeleteWishCommand) error {
	uof := uc.uofFactory.NewUnitOfWork(true)
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err := permissionService.CheckModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}
		wishRepo := rf.WishRepository()
		err = wishRepo.DeleteWish(ctx, cmd.WishlistID, cmd.WishID)
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

func (uc *WishUsecase) ReserveWish(ctx context.Context, cmd ReserveWishCommand) error {
	uof := uc.uofFactory.NewUnitOfWork(true)
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err := permissionService.CheckReservationInWishlist(
			ctx, cmd.RequestorUserID, cmd.WishlistID,
		)
		if err != nil {
			return err
		}
		wishRepo := rf.WishRepository()
		wish, err := wishRepo.GetWish(ctx, cmd.WishlistID, cmd.WishID)
		if err != nil {
			return err
		}
		reserveTime := time.Now()
		err = wish.Reserve(cmd.RequestorUserID, reserveTime)
		if err != nil {
			return err
		}
		err = wishRepo.UpdateWish(ctx, wish)
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
