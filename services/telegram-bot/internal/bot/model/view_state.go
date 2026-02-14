package model

import (
	"errors"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

var (
	ErrStateNotFound = errors.New("view state not found")
)

type State struct {
	ChatID int64
	Type   ui.StateType
}

func (s State) StateType() ui.StateType {
	return s.Type
}
