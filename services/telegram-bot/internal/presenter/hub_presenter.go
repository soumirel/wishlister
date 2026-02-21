package presenter

import (
	"context"
	"errors"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	ui_hub "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/hub"
)

type hubPresenter struct {
	navigator ui.Navigator
	view      ui_hub.View
}

func NewHubPresenter(
	navigator ui.Navigator,
) *hubPresenter {
	return &hubPresenter{
		navigator: navigator,
	}
}

func (p *hubPresenter) AttachView(view ui_hub.View) error {
	if p.view != nil {
		return errors.New("view already attached")
	}
	p.view = view
	return nil
}

func (p *hubPresenter) OnNavigateToWishlists(ctx context.Context, i ui_hub.NavigateToWishlistsIntent) error {
	err := p.navigator.NavigateToWishlists(ctx, ui.NavigateToWishlistsIntent{})
	if err != nil {
		return err
	}
	return nil
}

func (p *hubPresenter) OnNavigateToWishlistCreator(ctx context.Context, i ui_hub.NavigateToWishlistCreatorIntent) error {
	err := p.navigator.NavigateToWishlistCreator(ctx, ui.NavigateToWishlistCreatorIntent{})
	if err != nil {
		return err
	}
	return nil
}
