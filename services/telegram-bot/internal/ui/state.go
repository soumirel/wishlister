package ui

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrStateNotFound = errors.New("state not found")
)

type StateKey string

type StateStore interface {
	GetState(ctx context.Context, k StateKey) (json.RawMessage, error)
	StoreState(ctx context.Context, k StateKey, v json.RawMessage) error
	DeleteState(ctx context.Context, k StateKey) error
}
