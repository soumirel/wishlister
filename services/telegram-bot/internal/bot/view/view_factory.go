package view

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
)

type viewFactory struct {
	bot              *bot.Bot
	chatID           int64
	presenterFactory presenterFactory
	vlc              ui.ViewLifecycleController
}

func NewViewFactory(
	bot *bot.Bot,
	chatID int64,
	presenterFactory presenterFactory,
	vlc ui.ViewLifecycleController,
) *viewFactory {
	return &viewFactory{
		bot:              bot,
		chatID:           chatID,
		presenterFactory: presenterFactory,
		vlc:              vlc,
	}
}

func (f *viewFactory) NewHubView(ctx context.Context) (model.BotHubView, error) {
	v := NewHubView(
		f.bot,
		f.chatID,
		f.presenterFactory.HubPresenter(),
	)
	err := f.vlc.NewLifeCycle(ctx, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (f *viewFactory) NewWishlistsView(ctx context.Context) (model.BotWishlistsView, error) {
	v := NewWishlistsView(
		f.bot,
		f.chatID,
		f.presenterFactory.WishlistsPresenter(),
	)
	err := f.vlc.NewLifeCycle(ctx, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (f *viewFactory) NewWishlistCreatorView(ctx context.Context) (model.BotWishlistCreatorView, error) {
	v := NewWishlistCreatorView(
		f.bot,
		f.chatID,
		f.presenterFactory.WishlistCreatorPresenter(),
	)
	err := f.vlc.NewLifeCycle(ctx, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
