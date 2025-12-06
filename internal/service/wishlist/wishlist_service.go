package wishlist

import (
	"context"
	"errors"
	"log"
	"wishlister/internal/domain"
	"wishlister/internal/pkg"
)

// type reservationRepository interface {
// 	Reserve(ctx context.Context, reservation domain.Reservation) error
// }

// type wishlistShareRepository interface {
// 	CreateWishlistToken(ctx context.Context, token *domain.WishlistToken) error
// }

type wishlistService struct {
	wishlistRepo           wishlistRepository
	wishRepo               wishRepository
	wishlistPermissionRepo wishlistPermissionRepository
	transactor             pkg.CtxTransactor
	//reservationRepo reservationRepository
	//wishlistTokenRepo wishlistShareRepository
}

func NewWishlistService(
	wishlistRepo wishlistRepository,
	wishRepo wishRepository,
	wishlistPermissionRepo wishlistPermissionRepository,
	transactor pkg.CtxTransactor,
	//reservationRepo reservationRepository,
) *wishlistService {
	return &wishlistService{
		wishlistRepo:           wishlistRepo,
		wishRepo:               wishRepo,
		wishlistPermissionRepo: wishlistPermissionRepo,
		transactor:             transactor,
		//reservationRepo: reservationRepo,
	}
}

func (s *wishlistService) checkCanRead(ctx context.Context, userID, wishlistID string) (bool, error) {
	permission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx, userID, wishlistID)
	if err != nil {
		return false, err
	}
	if permission == nil {
		return false, nil
	}
	if !permission.Can(domain.ReadWishlistAction) {
		// TODO: Access denied error maybe?
		return false, nil
	}
	return true, nil
}

func (s *wishlistService) validateReadWishlistAccess(ctx context.Context, userID, wishlistID string) error {
	canRead, err := s.checkCanRead(ctx, userID, wishlistID)
	if err != nil {
		return err
	}
	if !canRead {
		// TODO: Access denied error maybe?
		return domain.ErrWishlistDoesNotExist
	}
	return nil
}

func (s *wishlistService) checkCanModify(ctx context.Context, userID, wishlistID string) (bool, error) {
	permission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx, userID, wishlistID)
	if err != nil {
		return false, err
	}
	if permission == nil {
		return false, nil
	}
	if !permission.Can(domain.ModifyWishlistAction) {
		// TODO: Access denied error maybe?
		return false, nil
	}
	return true, nil
}

func (s *wishlistService) validateModifyWishlistAccess(ctx context.Context, userID, wishlistID string) error {
	canModify, err := s.checkCanModify(ctx, userID, wishlistID)
	if err != nil {
		return err
	}
	if !canModify {
		// TODO: Access denied error maybe?
		return domain.ErrWishlistDoesNotExist
	}
	return nil
}

func (s *wishlistService) validateGrantPermissionCommand(ctx context.Context, cmd GrantWishlistPermissionCommand) error {
	if cmd.RequestorUserID == cmd.RequestingUserID {
		return errors.New("cannot create permission to yourself")
	}
	permission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx, cmd.RequestorUserID, cmd.WishlistID)
	if err != nil {
		return err
	}
	if permission == nil {
		return domain.ErrWishlistDoesNotExist
	}
	if !permission.CanGrantPermission(domain.WishlistPermissionLevel(cmd.PersmissionLevel)) {
		return errors.New("provided permission cannot be granted")
	}
	return nil
}

func (s *wishlistService) validateRevokePermissionCommand(ctx context.Context, cmd RevokeWishlistPermissionCommand) error {
	if cmd.RequestorUserID == cmd.RequestingUserID {
		return errors.New("can't delete permission from yourself")
	}
	// TODO: If more than one user can be owner for wishlist,
	// it is necessary to check for at least one owner
	requestorPermission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx,
		cmd.RequestorUserID, cmd.WishlistID,
	)
	if err != nil {
		return err
	}
	if requestorPermission == nil {
		return domain.ErrWishlistDoesNotExist
	}
	requestingPermission, err := s.wishlistPermissionRepo.GetPermissionToWishlist(ctx,
		cmd.RequestingUserID, cmd.WishlistID,
	)
	if err != nil {
		return err
	}
	if requestingPermission == nil {
		return domain.ErrWishlistPermissionNotExist
	}
	if !requestorPermission.CanRevokePermission(requestingPermission.Level) {
		return errors.New("provided permission cannot be revoked")
	}
	return nil
}

func (s *wishlistService) GetWishlists(ctx context.Context, cmd GetWishlistsCommand) ([]domain.Wishlist, error) {
	persmissions, err := s.wishlistPermissionRepo.GetUserPermissionsToWishlists(ctx, cmd.RequestorUserID)
	if err != nil {
		return nil, err
	}
	wishlistsIds := persmissions.GetWishlitsIdsForAction(domain.ReadWishlistAction)
	if len(wishlistsIds) == 0 {
		return nil, nil
	}
	wishlists, err := s.wishlistRepo.GetWishlists(ctx, WishlistFilter{
		WishlistIDs: wishlistsIds,
	})
	if err != nil {
		return nil, err
	}
	resp := make([]domain.Wishlist, 0, len(wishlists))
	for _, wl := range wishlists {
		resp = append(resp, *wl)
	}
	return resp, nil
}

func (s *wishlistService) GetWishlist(ctx context.Context, cmd GetWishlistCommand) (domain.Wishlist, error) {
	err := s.validateReadWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
	if err != nil {
		return domain.Wishlist{}, err
	}
	wishlist, err := s.wishlistRepo.GetWishlist(ctx, cmd.WishlistID)
	if err != nil {
		return domain.Wishlist{}, err
	}
	if wishlist == nil {
		return domain.Wishlist{}, domain.ErrWishDoesNotExist
	}
	return *wishlist, nil
}

func (s *wishlistService) CreateWishlist(ctx context.Context, cmd CreateWishlistCommand) (domain.Wishlist, error) {
	wishlist := domain.NewWishlist(cmd.RequestorUserID, cmd.Name)

	ctx, err := s.transactor.BeginCtxTx(ctx)
	if err != nil {
		return domain.Wishlist{}, nil
	}
	defer func() {
		if err == nil {
			return
		}
		err := s.transactor.RollbackCtxTx(ctx)
		if err != nil {
			log.Print(err.Error())
		}
	}()

	err = s.wishlistRepo.CreateWishlist(ctx, wishlist)
	if err != nil {
		return domain.Wishlist{}, err
	}
	persmission := domain.NewWishlistPersmission(
		wishlist.UserID, wishlist.ID, domain.OwnerWishlistPermissionLevel,
	)
	err = s.wishlistPermissionRepo.SaveWishlistPermission(ctx, persmission)
	if err != nil {
		return domain.Wishlist{}, err
	}

	s.transactor.CommitCtxTx(ctx)

	return *wishlist, nil
}

func (s *wishlistService) UpdateWishlist(ctx context.Context, cmd UpdateWishlistCommand) (domain.Wishlist, error) {
	err := s.validateModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
	if err != nil {
		return domain.Wishlist{}, err
	}
	// TODO: To DAO
	wishlist := &domain.Wishlist{
		ID:     cmd.WishlistID,
		UserID: cmd.RequestorUserID,
		Name:   cmd.Name,
	}
	err = s.wishlistRepo.UpdateWishlist(ctx, wishlist)
	if err != nil {
		return domain.Wishlist{}, err
	}
	return *wishlist, nil
}

func (s *wishlistService) DeleteWishlist(ctx context.Context, cmd DeleteWishlistCommand) error {
	err := s.validateModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
	if err != nil {
		return err
	}
	err = s.wishlistRepo.DeleteWishlist(ctx, cmd.WishlistID)
	if err != nil {
		return err
	}
	return nil
}

func (s *wishlistService) GrantWishlistPermission(ctx context.Context, cmd GrantWishlistPermissionCommand) error {
	err := s.validateGrantPermissionCommand(ctx, cmd)
	if err != nil {
		return err
	}
	permission := domain.NewWishlistPersmission(
		cmd.RequestingUserID, cmd.WishlistID, domain.WishlistPermissionLevel(cmd.PersmissionLevel),
	)
	err = s.wishlistPermissionRepo.SaveWishlistPermission(ctx, permission)
	if err != nil {
		return err
	}
	return nil
}

func (s *wishlistService) RevokeWishlistPermission(ctx context.Context, cmd RevokeWishlistPermissionCommand) error {
	err := s.validateRevokePermissionCommand(ctx, cmd)
	if err != nil {
		return err
	}
	err = s.wishlistPermissionRepo.DeleteWishlistPermission(ctx, cmd.RequestingUserID, cmd.WishlistID)
	if err != nil {
		return err
	}
	return nil
}

func (s *wishlistService) GetWish(ctx context.Context, cmd GetWishCommand) (domain.Wish, error) {
	err := s.validateReadWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
	if err != nil {
		return domain.Wish{}, err
	}
	wish, err := s.wishRepo.GetWish(ctx, cmd.RequestorUserID, cmd.WishlistID, cmd.WishID)
	if err != nil {
		return domain.Wish{}, err
	}
	if wish == nil {
		return domain.Wish{}, domain.ErrWishlistDoesNotExist
	}
	return *wish, nil
}

func (s *wishlistService) CreateWish(ctx context.Context, cmd CreateWishCommand) (domain.Wish, error) {
	err := s.validateModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
	if err != nil {
		return domain.Wish{}, err
	}
	wishlist, err := s.GetWishlist(ctx, GetWishlistCommand{
		RequestorUserID: cmd.RequestorUserID,
		WishlistID:      cmd.WishlistID,
	})
	if err != nil {
		return domain.Wish{}, err
	}
	wish, err := wishlist.AddWish(cmd.WishName)
	if err != nil {
		return domain.Wish{}, err
	}
	err = s.wishRepo.CreateWish(ctx, wish)
	if err != nil {
		return domain.Wish{}, err
	}
	return *wish, nil
}

func (s *wishlistService) UpdateWish(ctx context.Context, cmd UpdateWishCommand) (domain.Wish, error) {
	err := s.validateModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
	if err != nil {
		return domain.Wish{}, err
	}
	wishlist, err := s.GetWishlist(ctx, GetWishlistCommand{
		RequestorUserID: cmd.RequestorUserID,
		WishlistID:      cmd.WishlistID,
	})
	if err != nil {
		return domain.Wish{}, err
	}
	wish, err := wishlist.UpdateWish(cmd.WishID, cmd.WishName)
	if err != nil {
		return domain.Wish{}, err
	}
	err = s.wishRepo.UpdateWish(ctx, wish)
	if err != nil {
		return domain.Wish{}, err
	}
	return *wish, nil
}

func (s *wishlistService) DeleteWish(ctx context.Context, cmd DeleteWishCommand) error {
	err := s.validateModifyWishlistAccess(ctx, cmd.RequestorUserID, cmd.WishlistID)
	if err != nil {
		return err
	}
	wishlist, err := s.GetWishlist(ctx, GetWishlistCommand{
		RequestorUserID: cmd.RequestorUserID,
		WishlistID:      cmd.WishlistID,
	})
	if err != nil {
		return err
	}
	err = wishlist.DeleteWish(cmd.WishID)
	if err != nil {
		return err
	}
	err = s.wishRepo.DeleteWish(ctx, cmd.RequestorUserID, cmd.WishlistID, cmd.WishID)
	if err != nil {
		return err
	}
	return nil
}

// func (s *WishlistService) Reserve(ctx context.Context, reservation domain.Reservation) error {
// 	if reservation.UserID == "" || reservation.WishID == "" {
// 		return fmt.Errorf("bad id")
// 	}
// 	return s.reservationRepo.Reserve(ctx, reservation)
// }

// func (s *WishlistService) CreateShareToken(ctx context.Context) (domain.WishlistToken, error) {
// 	au, ok := auth.FromCtx(ctx)
// 	if !ok {
// 		return domain.WishlistToken{}, auth.ErrUnauthorized
// 	}
// 	token := domain.NewWishlistToken(au.UserID)
// 	err := s.wishlistTokenRepo.CreateWishlistToken(ctx, token)
// 	if err != nil {
// 		return domain.WishlistToken{}, err
// 	}
// 	return *token, nil
// }
