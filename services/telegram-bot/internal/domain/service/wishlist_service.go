package service

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/entity"
)

type GetWishlistsParams struct {
	RequestorUserID string
}

type WishlistCoreReadService interface {
	GetWishlists(ctx context.Context) (entity.WishlistListModel, error)
}
