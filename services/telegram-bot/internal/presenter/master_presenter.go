package presenter

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

type masterPresenter struct {
	svcFactory service.ServiceFactory
}

func NewMasterPresenter(svcFactory service.ServiceFactory) *masterPresenter {
	return &masterPresenter{
		svcFactory: svcFactory,
	}
}

func (p *masterPresenter) Module() ui.ModuleType {
	return ui.MasterModule
}

func (p *masterPresenter) HandleIntent(ctx context.Context, i ui.Intent) (ui.ViewModel, error) {
	switch ti := i.(type) {
	case ui.ShowWishlistsIntent:
		return p.handleShowWishlistsIntent(ctx, ti)
	default:
		return nil, ui.ErrUnknownIntent
	}
}

func (p *masterPresenter) handleShowWishlistsIntent(ctx context.Context, i ui.ShowWishlistsIntent) (ui.ViewModel, error) {
	wishlistSvc := p.svcFactory.GetWishlistCoreReadService()
	wishlists, err := wishlistSvc.GetWishlists(ctx)
	if err != nil {
		return nil, err
	}
	return ui.ShowWishlistsViewModel{
		Wishlists: wishlists,
	}, nil
}
