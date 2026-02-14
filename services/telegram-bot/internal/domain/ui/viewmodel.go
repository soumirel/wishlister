package ui

import (
	"errors"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
)

var (
	ErrUnknownViewModel = errors.New("unknown view model")
)

type ViewModel interface{}

type ShowWishlistsViewModel struct {
	Wishlists model.WishlistList
}

type CreateWishlistNameWaitingVM struct{}
