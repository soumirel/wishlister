package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	"github.com/valkey-io/valkey-go"
)

type vlsStore struct {
	valkeyClient valkey.Client
}

func NewVlsStore(
	valkeyClient valkey.Client,
) *vlsStore {
	return &vlsStore{
		valkeyClient: valkeyClient,
	}
}

func (s *vlsStore) buildKey(id ui.ViewIdentifier) string {
	return fmt.Sprintf("vls:%v", id)
}

func (s *vlsStore) GetViewLifecycleStatus(ctx context.Context, id ui.ViewIdentifier) (ui.ViewLifecycleStatus, error) {
	res, err := s.valkeyClient.Do(ctx,
		s.valkeyClient.B().Get().Key(s.buildKey(id)).Build(),
	).ToString()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return ui.VLS_Empty, nil
		}
		return ui.VLS_Empty, nil
	}
	return ui.ViewLifecycleStatus(res), nil
}

func (s *vlsStore) StoreViewLifecycleStatus(ctx context.Context, id ui.ViewIdentifier, vls ui.ViewLifecycleStatus) error {
	err := s.valkeyClient.Do(
		ctx, s.valkeyClient.B().Set().Key(s.buildKey(id)).Value(string(vls)).Build(),
	).Error()
	if err != nil {
		return err
	}
	return nil
}
