package service

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
)

type WishlistCoreReadService interface {
	GetWishlists(ctx context.Context) (model.WishlistList, error)
}

type CreateWishlistParams struct {
	Name string
}

type WishlistCoreQueryService interface {
	CreateWishlist(ctx context.Context, params CreateWishlistParams) (model.Wishlist, error)
}

type WishlistCoreService interface {
	WishlistCoreQueryService
	WishlistCoreReadService
}
