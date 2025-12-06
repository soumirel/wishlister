package wishlist

import (
	"context"
	"wishlister/internal/domain"
)

type WishlistFilter struct {
	WishlistIDs []string
}

type wishlistRepository interface {
	GetUserWishlists(ctx context.Context, userID string) ([]*domain.Wishlist, error)
	GetWishlists(ctx context.Context, filter WishlistFilter) ([]*domain.Wishlist, error)
	GetWishlist(ctx context.Context, wishlistID string) (*domain.Wishlist, error)
	CreateWishlist(ctx context.Context, wishlist *domain.Wishlist) error
	UpdateWishlist(ctx context.Context, wishlst *domain.Wishlist) error
	DeleteWishlist(ctx context.Context, wishlistID string) error
}

type wishRepository interface {
	GetWish(ctx context.Context, userID, wishlistID, wishID string) (*domain.Wish, error)
	CreateWish(ctx context.Context, wish *domain.Wish) error
	UpdateWish(ctx context.Context, wish *domain.Wish) error
	DeleteWish(ctx context.Context, userID, wishlistID, wishID string) error
}

type wishlistPermissionRepository interface {
	GetPermissionToWishlist(ctx context.Context, userID, wishlistID string) (*domain.WishlistPermission, error)
	GetUserPermissionsToWishlists(ctx context.Context, userID string) (domain.WishlistsPersmmissions, error)
	SaveWishlistPermission(context.Context, *domain.WishlistPermission) error
	DeleteWishlistPermission(ctx context.Context, userID, wishlistID string) error
}
