package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
	"github.com/valkey-io/valkey-go"
)

type stateRepository struct {
	valkeyClient valkey.Client
}

func NewViewStateRepository(
	valkeyClient valkey.Client,
) *stateRepository {
	return &stateRepository{
		valkeyClient: valkeyClient,
	}
}

func (r *stateRepository) SaveState(ctx context.Context, s model.State) error {
	err := r.valkeyClient.Do(ctx,
		r.valkeyClient.B().Set().Key(r.buildStateKeyFromState(s)).Value(string(s.Type)).Build(),
	).Error()
	if err != nil {
		return err
	}
	return nil
}

func (r *stateRepository) GetState(ctx context.Context, chatID int64) (model.State, error) {
	resp, err := r.valkeyClient.Do(ctx,
		r.valkeyClient.B().Get().Key(r.buildStateKey(chatID)).Build(),
	).ToString()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return model.State{}, model.ErrStateNotFound
		}
		return model.State{}, err
	}
	return model.State{
		ChatID: chatID,
		Type:   ui.StateType(resp),
	}, nil
}

func (r *stateRepository) buildStateKeyFromState(s model.State) string {
	return fmt.Sprintf("state:%v", s.ChatID)
}

func (r *stateRepository) buildStateKey(chatID int64) string {
	return fmt.Sprintf("state:%v", chatID)
}
