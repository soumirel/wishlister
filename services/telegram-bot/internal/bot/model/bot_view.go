package model

import (
	"context"

	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui/hub"
	wishlistcreator "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlist_creator"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlists"
)

type UpdateProcessor interface {
	ProcessUpdate(ctx context.Context, update *models.Update) error
}

type CommandProcessor interface {
	ProcessCommand(ctx context.Context, cmd Command) error
}

type BotHubView interface {
	hub.View
	UpdateProcessor
	CommandProcessor
}

type BotWishlistsView interface {
	wishlists.View
	UpdateProcessor
}

type BotWishlistCreatorView interface {
	wishlistcreator.View
	UpdateProcessor
}

type ViewRegistry interface {
	HubView(ctx context.Context) (BotHubView, error)
	WishlistsView(ctx context.Context) (BotWishlistsView, error)
	WishlistCreatorView(ctx context.Context) (BotWishlistCreatorView, error)
}

type ViewFactory interface {
	NewHubView(ctx context.Context) (BotHubView, error)
	NewWishlistsView(ctx context.Context) (BotWishlistsView, error)
	NewWishlistCreatorView(ctx context.Context) (BotWishlistCreatorView, error)
}
