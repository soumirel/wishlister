package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	"github.com/valkey-io/valkey-go"
)

type storeContainer struct {
	Data json.RawMessage `json:"data"`
	// Version   int64           `json:"_version"`
	// CreatedAt int64           `json:"_created"`
	// UpdatedAt int64           `json:"_updated"`
}

type stateStore struct {
	valkeyClient valkey.Client
}

func NewStateStore(
	valkeyClient valkey.Client,
) *stateStore {
	return &stateStore{
		valkeyClient: valkeyClient,
	}
}

func (s *stateStore) buildStoreKey(k ui.StateKey) string {
	return fmt.Sprintf("view:%v", k)
}

func (s *stateStore) GetState(ctx context.Context, k ui.StateKey) (json.RawMessage, error) {
	bytes, err := s.valkeyClient.Do(ctx,
		s.valkeyClient.B().Get().Key(s.buildStoreKey(k)).Build(),
	).AsBytes()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return nil, ui.ErrStateNotFound
		}
		return nil, err
	}
	var container storeContainer
	err = json.Unmarshal(bytes, &container)
	if err != nil {
		return nil, err
	}
	return container.Data, nil
}

func (s *stateStore) StoreState(ctx context.Context, k ui.StateKey, v json.RawMessage) error {
	containter := storeContainer{
		Data: v,
	}
	bytes, err := json.Marshal(containter)
	if err != nil {
		return err
	}
	err = s.valkeyClient.Do(
		ctx, s.valkeyClient.B().Set().Key(s.buildStoreKey(k)).Value(valkey.BinaryString(bytes)).Build(),
	).Error()
	if err != nil {
		return err
	}
	return nil
}

func (s *stateStore) DeleteState(ctx context.Context, k ui.StateKey) error {
	err := s.valkeyClient.Do(
		ctx, s.valkeyClient.B().Del().Key(s.buildStoreKey(k)).Build(),
	).Error()
	if err != nil {
		return err
	}
	return nil
}
