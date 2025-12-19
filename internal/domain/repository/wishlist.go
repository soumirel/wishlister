package repository

import (
	"context"
	entity "wishlister/internal/domain/entity"
)

type WishlistRepository interface {
	GetWishlists(ctx context.Context, wishlistsIDs []string) ([]*entity.Wishlist, error)
	GetWishlist(ctx context.Context, wishlistID string) (*entity.Wishlist, error)
	CreateWishlist(ctx context.Context, wishlist *entity.Wishlist) error
	UpdateWishlist(ctx context.Context, wishlst *entity.Wishlist) error
	DeleteWishlist(ctx context.Context, wishlistID string) error
}
