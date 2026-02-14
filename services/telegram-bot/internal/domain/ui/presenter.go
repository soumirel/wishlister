package ui

import (
	"context"
)

type PresenterType uint8

const (
	_ PresenterType = iota
	MasterPresenterType
)

type Presenter interface {
	HandleIntent(ctx context.Context, i Intent) (ViewModel, error)
	Module() ModuleType
}

type PresenterProvider interface {
	GetPresenter(m ModuleType) (Presenter, error)
}
