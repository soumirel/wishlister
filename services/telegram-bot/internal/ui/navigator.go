package ui

import (
	"context"
)

type NavigateToWishlistsIntent struct {
}

type NavigateToWishlistCreatorIntent struct {
}

type NavigateToHubIntent struct {
}

type Navigator interface {
	NavigateToHub(context.Context, NavigateToHubIntent) error
	NavigateToWishlists(context.Context, NavigateToWishlistsIntent) error
	NavigateToWishlistCreator(context.Context, NavigateToWishlistCreatorIntent) error
}
