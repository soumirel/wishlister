package telegrambot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
)

func StartTelegramBot(
	appCtx context.Context,
	botToken string,
	authSvc service.WishlistAuthService,
	wishlistReadSvc service.WishlistCoreReadService,
) error {
	mwFactory := newMiddlewareFactory(authSvc)

	botHandler := NewBotHandler(wishlistReadSvc)

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.Handle),
		bot.WithMiddlewares(
			mwFactory.AuthMiddleware(),
		),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		return err
	}

	b.RegisterHandler(
		bot.HandlerTypeMessageText, "wishlists",
		bot.MatchTypeCommandStartOnly, botHandler.HandleListCommand,
	)

	go b.Start(appCtx)

	return nil
}
