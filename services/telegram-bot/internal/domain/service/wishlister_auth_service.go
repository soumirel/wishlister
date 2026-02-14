package service

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
)

type WishlistAuthService interface {
	AuthByTelegramID(ctx context.Context, telegramID int64) (*model.WishlisterUser, error)
}
