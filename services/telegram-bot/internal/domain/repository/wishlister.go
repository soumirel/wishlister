package repository

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
)

type WishlistCoreAuthRepository interface {
	GetUserIdByExternalIdentity(ctx context.Context, ei model.ExternalIdentity) (string, error)
	CreateUserFromExternalIdentity(ctx context.Context, ei model.ExternalIdentity) (string, error)
}

type WishlistCoreRepository interface {
	WishlistCoreReadRepository
	WishlistCoreQueryRepository
}

type WishlistCoreReadRepository interface {
	GetWishlists(ctx context.Context) (model.WishlistList, error)
}

type WishlistCoreQueryRepository interface {
	CreateWishlist(ctx context.Context, w model.Wishlist) (model.Wishlist, error)
}
