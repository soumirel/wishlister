package service

import (
	"context"
	"errors"

	"github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/wishlist/internal/domain/repository"
	"github.com/soumirel/wishlister/wishlist/internal/domain/service"
)

type wishlistPermissionService struct {
	wishlistPermissionRepo repository.WishlistPermissionRepository
}

func newWishlistPersmissionService(
	wishlistPermissionRepo repository.WishlistPermissionRepository,
) *wishlistPermissionService {
	return &wishlistPermissionService{
		wishlistPermissionRepo: wishlistPermissionRepo,
	}
}

func (s *wishlistPermissionService) SaveWishlistPermission(ctx context.Context, p *entity.WishlistPermission) error {
	err := s.wishlistPermissionRepo.SaveWishlistPermission(ctx, p)
	if err != nil {
		return err
	}
	return err
}

func (s *wishlistPermissionService) GrantWishlistPermission(ctx context.Context, o service.GrantWishlistPermissionOptions) error {
	err := s.validatePermissionGranting(ctx, o)
	if err != nil {
		return err
	}
	permission := entity.NewWishlistPersmission(
		o.TargetUserID, o.WishlistID, entity.WishlistPermissionLevel(o.PermissionLevel),
	)
	err = s.wishlistPermissionRepo.SaveWishlistPermission(ctx, permission)
	if err != nil {
		return err
	}
	return nil
}

func (s *wishlistPermissionService) RevokeWishlistPermission(ctx context.Context, o service.RevokeWishlistPermissionOptions) error {
	err := s.validatePermissionRevoking(ctx, o)
	if err != nil {
		return err
	}
	err = s.wishlistPermissionRepo.DeleteWishlistPermission(ctx, o.TargetUserID, o.WishlistID)
	if err != nil {
		return err
	}
	return nil
}

func (s *wishlistPermissionService) CanReadWishlist(ctx context.Context, userID, wishlistID string) (bool, error) {
	permission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx, userID, wishlistID)
	if err != nil {
		return false, err
	}
	if permission == nil {
		return false, nil
	}
	if !permission.Can(entity.ReadWishlistAction) {
		// TODO: Access denied error maybe?
		return false, nil
	}
	return true, nil
}

func (s *wishlistPermissionService) CheckReadWishlistAccess(ctx context.Context, userID, wishlistID string) error {
	canRead, err := s.CanReadWishlist(ctx, userID, wishlistID)
	if err != nil {
		return err
	}
	if !canRead {
		// TODO: Access denied error maybe?
		return entity.ErrWishlistDoesNotExist
	}
	return nil
}

func (s *wishlistPermissionService) CanModifyWishlist(ctx context.Context, userID, wishlistID string) (bool, error) {
	permission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx, userID, wishlistID)
	if err != nil {
		return false, err
	}
	if permission == nil {
		return false, nil
	}
	if !permission.Can(entity.ModifyWishlistAction) {
		// TODO: Access denied error maybe?
		return false, nil
	}
	return true, nil
}

func (s *wishlistPermissionService) CheckModifyWishlistAccess(ctx context.Context, userID, wishlistID string) error {
	canModify, err := s.CanModifyWishlist(ctx, userID, wishlistID)
	if err != nil {
		return err
	}
	if !canModify {
		// TODO: Access denied error maybe?
		return entity.ErrWishlistDoesNotExist
	}
	return nil
}

func (s *wishlistPermissionService) CanReserveInWishlist(ctx context.Context, userID, wishlistID string) (bool, error) {
	permission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx, userID, wishlistID)
	if err != nil {
		return false, err
	}
	if permission == nil {
		return false, nil
	}
	if !permission.Can(entity.ReserveWishWishlistAction) {
		return false, nil
	}
	return true, nil
}

func (s *wishlistPermissionService) CheckReservationInWishlist(ctx context.Context, userID, wishlistID string) error {
	canReserve, err := s.CanReserveInWishlist(ctx, userID, wishlistID)
	if err != nil {
		return err
	}
	if !canReserve {
		return entity.ErrWishlistDoesNotExist
	}
	return nil
}

func (s *wishlistPermissionService) GetPermissionsToWishlists(ctx context.Context, userID string) (entity.WishlistsPermissions, error) {
	permissions, err := s.wishlistPermissionRepo.GetPermissionsToWishlists(ctx, userID)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *wishlistPermissionService) validatePermissionGranting(ctx context.Context, o service.GrantWishlistPermissionOptions) error {
	if o.RequestorUserID == o.TargetUserID {
		return errors.New("cannot create permission to yourself")
	}
	permission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx, o.RequestorUserID, o.WishlistID)
	if err != nil {
		return err
	}
	if permission == nil {
		return entity.ErrWishlistDoesNotExist
	}
	if !permission.CanGrantPermission(entity.WishlistPermissionLevel(o.PermissionLevel)) {
		return errors.New("provided permission cannot be granted")
	}
	return nil
}

func (s *wishlistPermissionService) validatePermissionRevoking(ctx context.Context, o service.RevokeWishlistPermissionOptions) error {
	if o.RequestorUserID == o.TargetUserID {
		return errors.New("can't delete permission from yourself")
	}
	// TODO: If more than one user can be owner for wishlist,
	// it is necessary to check for at least one owner
	requestorPermission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx,
		o.RequestorUserID, o.WishlistID,
	)
	if err != nil {
		return err
	}
	if requestorPermission == nil {
		return entity.ErrWishlistDoesNotExist
	}
	requestingPermission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx,
		o.TargetUserID, o.WishlistID,
	)
	if err != nil {
		return err
	}
	if requestingPermission == nil {
		return entity.ErrWishlistPermissionNotExist
	}
	if !requestorPermission.CanRevokePermission(requestingPermission.Level) {
		return errors.New("provided permission cannot be revoked")
	}
	return nil
}
