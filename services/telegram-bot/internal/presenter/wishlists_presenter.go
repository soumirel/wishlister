package presenter

import (
	"context"
	"errors"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	ui_wishlists "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlists"
)

type wishlistsPresenter struct {
	navigator           ui.Navigator
	view                ui_wishlists.View
	wishlistReadService service.WishlistCoreReadService
}

func NewWishlistsPresenter(
	navigator ui.Navigator,
	wishlistReadService service.WishlistCoreReadService,
) *wishlistsPresenter {
	return &wishlistsPresenter{
		navigator:           navigator,
		wishlistReadService: wishlistReadService,
	}
}

func (p *wishlistsPresenter) AttachView(view ui_wishlists.View) error {
	if p.view != nil {
		return errors.New("view already attached")
	}
	p.view = view
	return nil
}

func (p *wishlistsPresenter) OnStart(ctx context.Context) error {
	wishlists, err := p.wishlistReadService.GetWishlists(ctx)
	if err != nil {
		return err
	}
	err = p.view.DisplayWishlists(ctx, ui_wishlists.WishlistsVM{
		Wishlists: wishlists,
	})
	if err != nil {
		return err
	}
	err = p.navigator.NavigateToHub(ctx, ui.NavigateToHubIntent{})
	if err != nil {
		return err
	}
	return nil
}
