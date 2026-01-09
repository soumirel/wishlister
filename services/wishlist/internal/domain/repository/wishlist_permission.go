package repository

import (
	"context"

	entity "github.com/soumirel/wishlister/wishlist/internal/domain/entity"
)

type WishlistPermissionRepository interface {
	GetPermissionToWishlist(ctx context.Context, userID, wishlistID string) (*entity.WishlistPermission, error)
	GetPermissionsToWishlists(ctx context.Context, userID string) (entity.WishlistsPermissions, error)
	SaveWishlistPermission(context.Context, *entity.WishlistPermission) error
	DeleteWishlistPermission(ctx context.Context, userID, wishlistID string) error
}
