package view

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
)

type viewRegistry struct {
	bot              *bot.Bot
	chatID           int64
	presenterFactory presenterFactory
	vlc              ui.ViewLifecycleController
}

func NewViewRegistry(
	bot *bot.Bot,
	chatID int64,
	presenterFactory presenterFactory,
	vlc ui.ViewLifecycleController,
) *viewRegistry {
	return &viewRegistry{
		bot:              bot,
		chatID:           chatID,
		presenterFactory: presenterFactory,
		vlc:              vlc,
	}
}

func (f *viewRegistry) HubView(ctx context.Context) (model.BotHubView, error) {
	v := NewHubView(
		f.bot,
		f.chatID,
		f.presenterFactory.HubPresenter(),
	)
	err := f.vlc.ContinueLifecycle(ctx, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (f *viewRegistry) WishlistsView(ctx context.Context) (model.BotWishlistsView, error) {
	v := NewWishlistsView(
		f.bot,
		f.chatID,
		f.presenterFactory.WishlistsPresenter(),
	)
	err := f.vlc.ContinueLifecycle(ctx, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (f *viewRegistry) WishlistCreatorView(ctx context.Context) (model.BotWishlistCreatorView, error) {
	v := NewWishlistCreatorView(
		f.bot,
		f.chatID,
		f.presenterFactory.WishlistCreatorPresenter(),
	)
	err := f.vlc.ContinueLifecycle(ctx, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
