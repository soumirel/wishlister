package service

import (
	"context"

	"github.com/soumirel/wishlister/telegram-bot/internal/domain/entity"
)

type WishlisterAuthService interface {
	AuthByTelegramID(ctx context.Context, telegramID int64) (*entity.WishlisterUser, error)
}
