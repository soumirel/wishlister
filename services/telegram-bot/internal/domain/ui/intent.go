package ui

import (
	"context"
)

type Intent interface{}

type IntentDispatcher interface {
	DispatchViewIntent(ctx context.Context, i Intent, v View) error
}
