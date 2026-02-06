package service

import (
	"context"

	"github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
)

type GrantWishlistPermissionOptions struct {
	RequestorUserID string
	WishlistID      string
	TargetUserID    string
	PermissionLevel string
}

type RevokeWishlistPermissionOptions struct {
	RequestorUserID string
	WishlistID      string
	TargetUserID    string
}

type WishlistPermissionService interface {
	SaveWishlistPermission(ctx context.Context, p *entity.WishlistPermission) error

	GrantWishlistPermission(ctx context.Context, o GrantWishlistPermissionOptions) error
	RevokeWishlistPermission(ctx context.Context, o RevokeWishlistPermissionOptions) error

	CanReadWishlist(ctx context.Context, userID, wishlistID string) (bool, error)
	CanModifyWishlist(ctx context.Context, userID, wishlistID string) (bool, error)
	CanReserveInWishlist(ctx context.Context, userID, wishlistID string) (bool, error)

	CheckReadWishlistAccess(ctx context.Context, userID, wishlistID string) error
	CheckModifyWishlistAccess(ctx context.Context, userID, wishlistID string) error

	CheckReservationInWishlist(ctx context.Context, userID, wishlistID string) error

	GetPermissionsToWishlists(ctx context.Context, userID string) (entity.WishlistsPermissions, error)
}
