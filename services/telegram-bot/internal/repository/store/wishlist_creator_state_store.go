package store

import (
	"context"
	"encoding/json"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	ui_wishlistcreator "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlist_creator"
)

type wishlistCreatorStateStore struct {
	stateStore ui.StateStore
}

func NewWishlistCreatorStateStore(stateStore ui.StateStore) *wishlistCreatorStateStore {
	return &wishlistCreatorStateStore{
		stateStore: stateStore,
	}
}

func (s *wishlistCreatorStateStore) GetWishlistCreatorState(ctx context.Context, k ui.StateKey) (ui_wishlistcreator.State, error) {
	raw, err := s.stateStore.GetState(ctx, k)
	if err != nil {
		return ui_wishlistcreator.State{}, err
	}
	var state ui_wishlistcreator.State
	err = json.Unmarshal(raw, &state)
	if err != nil {
		return state, err
	}
	return state, nil
}

func (s *wishlistCreatorStateStore) StoreWishlistCreatorState(ctx context.Context, k ui.StateKey, state ui_wishlistcreator.State) error {
	raw, err := json.Marshal(state)
	if err != nil {
		return nil
	}
	err = s.stateStore.StoreState(ctx, k, json.RawMessage(raw))
	if err != nil {
		return err
	}
	return nil
}

func (s *wishlistCreatorStateStore) DeleteWishlistCreatorState(ctx context.Context, k ui.StateKey) error {
	err := s.stateStore.DeleteState(ctx, k)
	if err != nil {
		return err
	}
	return nil
}
