package service

import (
	"context"

	"github.com/soumirel/wishlister/telegram-bot/internal/domain/entity"
)

type GetWishlistsParams struct {
	RequestorUserID string
}

type WishlisterReadService interface {
	GetWishlists(ctx context.Context) ([]*entity.Wishlist, error)
}
