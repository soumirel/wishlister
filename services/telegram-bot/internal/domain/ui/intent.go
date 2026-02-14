package ui

import (
	"context"
	"errors"
)

var (
	ErrUnknownIntent = errors.New("unknown intent")
)

type Intent interface{}

type IntentDispatcher interface {
	DispatchViewIntent(ctx context.Context, i Intent, v View) error
}

type ShowWishlistsIntent struct {
}

type CreateWishlistIntent struct {
}
