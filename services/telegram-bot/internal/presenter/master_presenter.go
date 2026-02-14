package presenter

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

type masterPresenter struct {
	svcFactory service.ServiceFactory
}

/*
1. Получили Intent
Чтобы понять, что нам делать, нужно понять стейт экрана


*/

func NewMasterPresenter(svcFactory service.ServiceFactory) *masterPresenter {
	return &masterPresenter{
		svcFactory: svcFactory,
	}
}

func (p *masterPresenter) Module() ui.ModuleType {
	return ui.MasterModule
}

func (p *masterPresenter) HandleIntent(ctx context.Context, i ui.Intent, st ui.State) (ui.ViewModel, error) {
	switch ti := i.(type) {
	case ui.ShowWishlistsIntent:
		return p.handleShowWishlistsIntent(ctx, ti)
	case ui.CreateWishlistIntent:
		return p.handleCreateWishlistIntent(ctx, ti)
	default:
		return nil, ui.ErrUnknownIntent
	}
}

func (p *masterPresenter) handleShowWishlistsIntent(ctx context.Context, _ ui.ShowWishlistsIntent) (ui.ViewModel, error) {
	wishlistSvc := p.svcFactory.GetWishlistCoreReadService()
	wishlists, err := wishlistSvc.GetWishlists(ctx)
	if err != nil {
		return nil, err
	}
	return ui.ShowWishlistsViewModel{
		Wishlists: wishlists,
	}, nil
}

func (p *masterPresenter) handleCreateWishlistIntent(ctx context.Context, i ui.CreateWishlistIntent) (ui.ViewModel, error) {
	return ui.CreateWishlistNameWaitingVM{}, nil
}
