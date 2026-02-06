package wishlist_permission

import (
	"context"

	"github.com/soumirel/wishlister/services/wishlist/internal/domain/repository"
	svcentity "github.com/soumirel/wishlister/services/wishlist/internal/domain/service"
	"github.com/soumirel/wishlister/services/wishlist/internal/service"
	"github.com/soumirel/wishlister/services/wishlist/internal/usecase"
)

type WishlistPermissionUsecase struct {
	uofFactory usecase.UnitOfWorkFactory
}

func NewWishlistPermissionUsecase(
	uofFactory usecase.UnitOfWorkFactory,
) *WishlistPermissionUsecase {
	return &WishlistPermissionUsecase{
		uofFactory: uofFactory,
	}
}

func (uc *WishlistPermissionUsecase) GrantWishlistPermission(ctx context.Context, cmd GrantWishlistPermissionCommand) error {
	uof := uc.uofFactory.NewUnitOfWork(true)
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()

		err := permissionService.CheckModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}

		err = permissionService.GrantWishlistPermission(ctx, svcentity.GrantWishlistPermissionOptions{
			RequestorUserID: cmd.RequestorUserID,
			WishlistID:      cmd.WishlistID,
			TargetUserID:    cmd.TargetUserID,
			PermissionLevel: cmd.PermissionLevel,
		})
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

func (uc *WishlistPermissionUsecase) RevokeWishlistPermissionCommand(ctx context.Context, cmd RevokeWishlistPermissionCommand) error {
	uof := uc.uofFactory.NewUnitOfWork(true)
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		permissionService := svcFactory.WishlistPermissionService()
		err := permissionService.CheckModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
		if err != nil {
			return err
		}

		err = permissionService.RevokeWishlistPermission(ctx, svcentity.RevokeWishlistPermissionOptions{
			RequestorUserID: cmd.RequestorUserID,
			WishlistID:      cmd.WishlistID,
			TargetUserID:    cmd.TargetUserID,
		})
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
