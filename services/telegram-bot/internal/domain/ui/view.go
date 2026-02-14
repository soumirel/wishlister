package ui

import (
	"context"
)

type View interface {
	Display(ctx context.Context, vm ViewModel) error
	State() State
	Module() ModuleType
}
