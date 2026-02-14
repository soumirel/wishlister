package service

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
)

type GetWishlistsParams struct {
	RequestorUserID string
}

type WishlistCoreReadService interface {
	GetWishlists(ctx context.Context) (model.WishlistList, error)
}
