package handler

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/auth"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
)

type middlewareFactory struct {
	authSvc service.WishlistAuthService
}

func newMiddlewareFactory(
	authSvc service.WishlistAuthService,
) *middlewareFactory {
	return &middlewareFactory{
		authSvc: authSvc,
	}
}

func (mwf *middlewareFactory) AuthMiddleware() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			user, err := mwf.authUpdate(ctx, update)
			if err != nil {
				log.Print(fmt.Errorf("auth error: %w", err))
				return
			}
			if user == nil {
				log.Print("no wishlister user, abort handling")
				return
			}
			ctx = auth.NewCtx(ctx, auth.Auth{
				UserID: user.ID,
			})
			next(ctx, b, update)
		}
	}
}

func (mwf *middlewareFactory) authUpdate(ctx context.Context, update *models.Update) (*model.WishlisterUser, error) {
	senderTelegramUserID, ok := getTelegramSenderUserID(update)
	if !ok {
		return nil, errors.New("no sender id in update")
	}
	wishlisterUser, err := mwf.authSvc.AuthByTelegramID(ctx, senderTelegramUserID)
	if err != nil {
		return nil, err
	}
	return wishlisterUser, nil
}

func getTelegramSenderUserID(update *models.Update) (int64, bool) {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return 0, false
	}
	return update.Message.From.ID, true
}
