package telegrambot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
)

func StartTelegramBot(
	appCtx context.Context,
	botToken string,
	authSvc service.WishlisterAuthService,
) error {
	mwFactory := newMiddlewareFactory(authSvc)

	botHandler := NewBotHandler()

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

	go b.Start(appCtx)

	return nil
}
