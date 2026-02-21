package model

import (
	"context"
	"errors"
)

type ViewType string

const (
	HubViewType             = "hub"
	WishlistsViewType       = "wishlists"
	WishlistCreatorViewType = "wishlistcreator"
)

const (
	tgbotIdPrefix = "tgbot"

	HubIdPrefix             = tgbotIdPrefix + ":hub"
	WishlistsIdPrefix       = tgbotIdPrefix + ":wishlists"
	WishlistCreatorIdPrefix = tgbotIdPrefix + ":wishlistcreator"
)

var (
	ErrActiveViewNotFound = errors.New("active view not found")
)

type ViewStore interface {
	GetFocusedViewType(ctx context.Context, chatID int64) (ViewType, error)
	StoreFocusedViewType(ctx context.Context, chatID int64, vt ViewType) error
}
