package presenter

import (
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui/hub"
	ui_wishlistcreator "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlist_creator"
	wishlistcreator "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlist_creator"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlists"
)

type presenterFactory struct {
	navigator       ui.Navigator
	wishlistService service.WishlistCoreService
	stateQueryStore ui_wishlistcreator.StateQueryStore
}

func NewPresenterFactory(
	navigator ui.Navigator,
	wishlistService service.WishlistCoreService,
	stateQueryStore ui_wishlistcreator.StateQueryStore,
) *presenterFactory {
	return &presenterFactory{
		navigator:       navigator,
		wishlistService: wishlistService,
		stateQueryStore: stateQueryStore,
	}
}

func (f *presenterFactory) HubPresenter() hub.Presenter {
	return NewHubPresenter(
		f.navigator,
	)
}

func (f *presenterFactory) WishlistsPresenter() wishlists.Presenter {
	return NewWishlistsPresenter(
		f.navigator,
		f.wishlistService,
	)
}

func (f *presenterFactory) WishlistCreatorPresenter() wishlistcreator.Presenter {
	return NewWishlistCreationPresenter(
		f.navigator,
		f.stateQueryStore,
		f.wishlistService,
	)
}
