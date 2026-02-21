package handler

import (
	"context"
	"errors"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/gofrs/uuid/v5"
	"github.com/soumirel/go-telegram-null-safety/wrap"
	"github.com/soumirel/wishlister/pkg/logger"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/auth"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
	"go.uber.org/zap"
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

func (mwf *middlewareFactory) LoggerMiddleware() bot.Middleware {
	mwLogger := logger.L().
		Named("bot")
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			requestID := uuid.Must(uuid.NewV7())
			reqLogger := mwLogger.With(
				zap.Int64("chat_id", wrap.Update(update).GetMessage().GetChat().GetID()),
				zap.String("request_id", requestID.String()),
			)
			ctx = logger.WithContext(ctx, reqLogger)
			next(ctx, bot, update)
			reqLogger.Info("bot_update")
		}
	}
}

func (mwf *middlewareFactory) ErrorHandler() bot.ErrorsHandler {
	return func(err error) {
		logger.L().Named("bot_server").Error(
			"bot server got error",
			zap.Error(err),
		)
	}
}

func (mwf *middlewareFactory) AuthMiddleware() bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			logger := logger.FromContext(ctx).Named("auth")
			user, err := mwf.authUpdate(ctx, update)
			if err != nil {
				logger.Warn("authentification error", zap.Error(err))
				return
			}
			if user == nil {
				logger.Warn("wishlister user is nil, abort handling")
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
