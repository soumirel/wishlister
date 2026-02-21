package view

import (
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui/hub"
	wishlistcreator "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlist_creator"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlists"
)

type presenterFactory interface {
	HubPresenter() hub.Presenter
	WishlistsPresenter() wishlists.Presenter
	WishlistCreatorPresenter() wishlistcreator.Presenter
}
