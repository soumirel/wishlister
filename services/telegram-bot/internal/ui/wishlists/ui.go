package wishlists

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
)

type WishlistsVM struct {
	Wishlists model.WishlistList
}

type State struct {
}

type Presenter interface {
	AttachView(View) error
	OnStart(context.Context) error
}

type View interface {
	ui.View
	DisplayWishlists(context.Context, WishlistsVM) error
}
