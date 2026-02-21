package navigator

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
)

type navigator struct {
	bot    *bot.Bot
	chatID int64

	viewFactory model.ViewFactory
	viewStore   model.ViewStore
}

func NewNavigator(
	bot *bot.Bot,
	chatID int64,
	viewStore model.ViewStore,
) *navigator {
	return &navigator{
		bot:       bot,
		chatID:    chatID,
		viewStore: viewStore,
	}
}

func (n *navigator) SetViewFactory(vf model.ViewFactory) {
	n.viewFactory = vf
}

func (n *navigator) NavigateToHub(ctx context.Context, i ui.NavigateToHubIntent) error {
	_, err := n.viewFactory.NewHubView(ctx)
	if err != nil {
		return err
	}
	err = n.viewStore.StoreFocusedViewType(ctx, n.chatID, model.HubViewType)
	if err != nil {
		return err
	}
	return nil
}

func (n *navigator) NavigateToWishlists(ctx context.Context, i ui.NavigateToWishlistsIntent) error {
	_, err := n.viewFactory.NewWishlistsView(ctx)
	if err != nil {
		return err
	}
	err = n.viewStore.StoreFocusedViewType(ctx, n.chatID, model.WishlistsViewType)
	if err != nil {
		return err
	}
	return nil
}

func (n *navigator) NavigateToWishlistCreator(ctx context.Context, i ui.NavigateToWishlistCreatorIntent) error {
	_, err := n.viewFactory.NewWishlistCreatorView(ctx)
	if err != nil {
		return err
	}
	err = n.viewStore.StoreFocusedViewType(ctx, n.chatID, model.WishlistCreatorViewType)
	if err != nil {
		return err
	}
	return nil
}
