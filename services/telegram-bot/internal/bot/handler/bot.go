package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

const (
	emptyPattern = ""
)

func StartTelegramBot(
	appCtx context.Context,
	botToken string,
	authSvc service.WishlistAuthService,
	intentDipatcher ui.IntentDispatcher,
	viewStateSvc viewStateService,
) error {
	mwFactory := newMiddlewareFactory(authSvc)
	handlerFactory := NewHandlerFactory(intentDipatcher, viewStateSvc)

	opts := []bot.Option{
		bot.WithDefaultHandler(
			handlerFactory.getEchoHandleFunc(),
		),
		bot.WithMiddlewares(
			mwFactory.AuthMiddleware(),
		),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		return err
	}

	handlerFactory.RegisterHandler(
		b,
		bot.HandlerTypeMessageText,
		wishlistsCommand,
		bot.MatchTypeCommandStartOnly,
	)

	handlerFactory.RegisterHandler(
		b,
		bot.HandlerTypeMessageText,
		createWishlistCommand,
		bot.MatchTypeCommandStartOnly,
	)

	go b.Start(appCtx)

	return nil
}
