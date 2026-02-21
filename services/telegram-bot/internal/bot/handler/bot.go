package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	wishlistcreator "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlist_creator"
)

const (
	noPattern = ""
)

func StartTelegramBot(
	appCtx context.Context,
	botToken string,
	authSvc service.WishlistAuthService,
	viewActivityStore model.ViewStore,
	wishlistService service.WishlistCoreService,
	stateQueryStore wishlistcreator.StateQueryStore,
	vlc ui.ViewLifecycleController,
) error {
	mwFactory := newMiddlewareFactory(authSvc)

	handlerFactory := NewBotHandleFuncFactory(
		viewActivityStore,
		wishlistService,
		stateQueryStore,
		vlc,
	)

	opts := []bot.Option{
		bot.WithDefaultHandler(
			handlerFactory.NewBotHandlerFunc(),
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
		model.WishlistsCommand,
		bot.MatchTypeCommandStartOnly,
	)

	handlerFactory.RegisterHandler(
		b,
		bot.HandlerTypeMessageText,
		model.CreateWishlistCommand,
		bot.MatchTypeCommandStartOnly,
	)

	go b.Start(appCtx)

	return nil
}
