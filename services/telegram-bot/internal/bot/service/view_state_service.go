package service

import (
	"context"
	"errors"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

type viewStateRepository interface {
	GetState(ctx context.Context, chatID int64) (model.State, error)
	SaveState(ctx context.Context, s model.State) error
}

type viewStateService struct {
	stateRepo viewStateRepository
}

func NewViewStateService(
	stateRepo viewStateRepository,
) *viewStateService {
	return &viewStateService{
		stateRepo: stateRepo,
	}
}

func (s *viewStateService) GetState(ctx context.Context, chatID int64) (model.State, error) {
	state, err := s.stateRepo.GetState(ctx, chatID)
	if err == nil {
		return state, nil
	}
	if errors.Is(err, model.ErrStateNotFound) {
		state = model.State{
			ChatID: chatID,
			Type:   ui.IdleStateType,
		}
		return state, s.SaveState(ctx, state)
	}
	return model.State{}, err
}
func (s *viewStateService) SaveState(ctx context.Context, st model.State) error {
	err := s.stateRepo.SaveState(ctx, st)
	if err != nil {
		return err
	}
	return nil
}
