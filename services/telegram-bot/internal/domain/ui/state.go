package ui

import "context"

type StateType string

const (
	IdleStateType StateType = "idle"

	// wishlist creation
	WishlistCreationNameWaiting StateType = "wishlist-creation:name-waiting"
)

type State interface {
	StateType() StateType
}

type StateStore interface {
	GetState(ctx context.Context) (State, error)
	StoreState(ctx context.Context, state State) error
}

type IdleState struct{}

type WishlistCreationState struct {
	Type StateType
}
