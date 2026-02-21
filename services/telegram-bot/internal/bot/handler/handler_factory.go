package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	tgmodels "github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/pkg/logger"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/navigator"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/view"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	wishlistcreator "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlist_creator"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/presenter"
)

type handlerFactory struct {
	viewActivityStore model.ViewStore
	wishlistService   service.WishlistCoreService
	stateQueryStore   wishlistcreator.StateQueryStore
	vlc               ui.ViewLifecycleController
}

func NewBotHandleFuncFactory(
	viewActivityStore model.ViewStore,
	wishlistService service.WishlistCoreService,
	stateQueryStore wishlistcreator.StateQueryStore,
	vlc ui.ViewLifecycleController,
) *handlerFactory {
	return &handlerFactory{
		viewActivityStore: viewActivityStore,
		wishlistService:   wishlistService,
		stateQueryStore:   stateQueryStore,
		vlc:               vlc,
	}
}

func (hf *handlerFactory) RegisterHandler(
	b *bot.Bot,
	handlerType bot.HandlerType,
	commandPattern model.Command,
	matchType bot.MatchType,
	m ...bot.Middleware,
) string {
	var (
		handlerFunc bot.HandlerFunc
	)

	if commandPattern == noPattern {
		handlerFunc = hf.NewBotHandlerFunc()
	} else if handlerType == bot.HandlerTypeCallbackQueryData {
		handlerFunc = hf.NewBotCallbackQueryHandlerFunc()
	} else if matchType == bot.MatchTypeCommand ||
		matchType == bot.MatchTypeCommandStartOnly {
		handlerFunc = hf.NewBotCommandHandlerFunc(commandPattern)
	} else {
		handlerFunc = hf.NewBotHandlerFunc()
	}

	return b.RegisterHandler(
		handlerType,
		string(commandPattern),
		matchType,
		handlerFunc,
		m...,
	)
}

func (hf *handlerFactory) NewBotHandlerFunc() bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgmodels.Update) {
		hf.botHandler(bot, update).HandleUpdate(ctx, bot, update)
	}
}

func (hf *handlerFactory) NewBotCallbackQueryHandlerFunc() bot.HandlerFunc {
	// TODO
	return func(ctx context.Context, bot *bot.Bot, update *tgmodels.Update) {
		logger.FromContext(ctx).Warn("got not currently supported callback query")
	}
}

func (hf *handlerFactory) NewBotCommandHandlerFunc(commandPattern model.Command) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgmodels.Update) {
		hf.botHandler(bot, update).HandleCommand(ctx, bot, update, commandPattern)
	}
}

func (hf *handlerFactory) botHandler(bot *bot.Bot, update *models.Update) *botHandler {
	navigator := navigator.NewNavigator(
		bot,
		update.Message.Chat.ID,
		hf.viewActivityStore,
	)
	presenterFactory := presenter.NewPresenterFactory(
		navigator,
		hf.wishlistService,
		hf.stateQueryStore,
	)
	viewRegistry := view.NewViewRegistry(
		bot,
		update.Message.Chat.ID,
		presenterFactory,
		hf.vlc,
	)
	viewFactory := view.NewViewFactory(
		bot,
		update.Message.Chat.ID,
		presenterFactory,
		hf.vlc,
	)
	navigator.SetViewFactory(viewFactory)
	botHandler := NewBotHandler(hf.viewActivityStore, viewRegistry)
	return botHandler
}
