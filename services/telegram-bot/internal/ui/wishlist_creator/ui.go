package wishlistcreator

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
)

type NameEnteringIntent struct {
	Name string
}

type NamePromtVM struct {
}

type CreatedVM struct {
	WishlistVM model.Wishlist
}

type State struct {
	Stage string `json:"stage"`
}

type StateReadStore interface {
	GetWishlistCreatorState(ctx context.Context, k ui.StateKey) (State, error)
}

type StateQueryStore interface {
	StoreWishlistCreatorState(ctx context.Context, k ui.StateKey, state State) error
	DeleteWishlistCreatorState(ctx context.Context, k ui.StateKey) error
}

type Presenter interface {
	AttachView(View) error
	OnStart(context.Context) error
	OnNameEntered(context.Context, NameEnteringIntent) error
	OnCancel(context.Context) error
}

type View interface {
	ui.View
	DisplayNamePromt(context.Context, NamePromtVM) error
	DisplayWishlistCreated(context.Context, CreatedVM) error
}
