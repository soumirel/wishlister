package ui

import (
	"context"
	"errors"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
)

var (
	ErrUnknownIntent = errors.New("unknown intent")

	ErrUnknownViewModel = errors.New("unknown view model")
)

type ViewModel interface{}

type View interface {
	Display(ctx context.Context, vm ViewModel) error
	Module() ModuleType
}

type ShowWishlistsIntent struct {
}

type ShowWishlistsViewModel struct {
	Wishlists model.WishlistList
}
