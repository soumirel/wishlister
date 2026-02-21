package presenter

import (
	"context"
	"errors"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	ui_wishlistcreator "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlist_creator"
)

const (
	wishlistCreatorWaitingForNameStage = "waiting_for_name"
)

type wishlistCreatorPresenter struct {
	navigator ui.Navigator
	view      ui_wishlistcreator.View

	stateQueryStore      ui_wishlistcreator.StateQueryStore
	wishlistQueryService service.WishlistCoreQueryService
}

func NewWishlistCreationPresenter(
	navigator ui.Navigator,
	stateQueryStore ui_wishlistcreator.StateQueryStore,
	wishlistQueryService service.WishlistCoreQueryService,
) *wishlistCreatorPresenter {
	return &wishlistCreatorPresenter{
		navigator:            navigator,
		stateQueryStore:      stateQueryStore,
		wishlistQueryService: wishlistQueryService,
	}
}

func (p *wishlistCreatorPresenter) AttachView(view ui_wishlistcreator.View) error {
	if p.view != nil {
		return errors.New("view already attached")
	}
	p.view = view
	return nil
}

func (p *wishlistCreatorPresenter) OnStart(ctx context.Context) error {
	state := ui_wishlistcreator.State{
		Stage: wishlistCreatorWaitingForNameStage,
	}
	err := p.view.DisplayNamePromt(ctx, ui_wishlistcreator.NamePromtVM{})
	if err != nil {
		return err
	}
	err = p.stateQueryStore.StoreWishlistCreatorState(ctx, ui.StateKey(p.view.ViewIdentifier()), state)
	if err != nil {
		return err
	}
	return nil
}

func (p *wishlistCreatorPresenter) OnNameEntered(ctx context.Context, i ui_wishlistcreator.NameEnteringIntent) error {
	wishlist, err := p.wishlistQueryService.CreateWishlist(ctx, service.CreateWishlistParams{
		Name: i.Name,
	})
	if err != nil {
		return err
	}
	err = p.view.DisplayWishlistCreated(ctx, ui_wishlistcreator.CreatedVM{
		WishlistVM: wishlist,
	})
	if err != nil {
		return err
	}
	err = p.stateQueryStore.DeleteWishlistCreatorState(ctx, ui.StateKey(p.view.ViewIdentifier()))
	if err != nil {
		return err
	}
	err = p.navigator.NavigateToHub(ctx, ui.NavigateToHubIntent{})
	if err != nil {
		return err
	}
	return nil
}

func (p *wishlistCreatorPresenter) OnCancel(ctx context.Context) error {
	err := p.stateQueryStore.DeleteWishlistCreatorState(ctx, ui.StateKey(p.view.ViewIdentifier()))
	if err != nil {
		return err
	}
	err = p.navigator.NavigateToHub(ctx, ui.NavigateToHubIntent{})
	if err != nil {
		return err
	}
	return nil
}
