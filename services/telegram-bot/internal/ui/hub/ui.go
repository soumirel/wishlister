package hub

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
)

type NavigateToWishlistsIntent struct {
}

type NavigateToWishlistCreatorIntent struct {
}

type Presenter interface {
	AttachView(View) error
	OnNavigateToWishlists(context.Context, NavigateToWishlistsIntent) error
	OnNavigateToWishlistCreator(context.Context, NavigateToWishlistCreatorIntent) error
}

type View interface {
	ui.View
}
