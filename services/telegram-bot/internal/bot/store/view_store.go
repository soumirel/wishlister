package store

import (
	"context"
	"fmt"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/valkey-io/valkey-go"
)

type viewStore struct {
	valkeyClient valkey.Client
}

func NewViewStore(
	valkeyClient valkey.Client,
) *viewStore {
	return &viewStore{
		valkeyClient: valkeyClient,
	}
}

func (v *viewStore) getKey(chatId int64) string {
	return fmt.Sprintf("%v", chatId)
}

func (v *viewStore) GetFocusedViewType(ctx context.Context, chatID int64) (model.ViewType, error) {
	resp, err := v.valkeyClient.Do(
		ctx, v.valkeyClient.B().Get().Key(v.getKey(chatID)).Build(),
	).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "", model.ErrActiveViewNotFound
		}
		return "", err
	}
	return model.ViewType(resp), nil
}

func (v *viewStore) StoreFocusedViewType(ctx context.Context, chatID int64, vt model.ViewType) error {
	err := v.valkeyClient.Do(
		ctx, v.valkeyClient.B().Set().Key(v.getKey(chatID)).Value(string(vt)).Build(),
	).Error()
	if err != nil {
		return err
	}
	return nil
}
