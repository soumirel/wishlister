package presenter

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

type intentDispatcher struct {
	presenterProvider ui.PresenterProvider
}

func NewIntentDispatcher(
	presenterProvider ui.PresenterProvider,
) *intentDispatcher {
	return &intentDispatcher{
		presenterProvider: presenterProvider,
	}
}

func (d *intentDispatcher) DispatchViewIntent(ctx context.Context, i ui.Intent, v ui.View) error {
	presenter, err := d.presenterProvider.GetPresenter(v.Module())
	if err != nil {
		return err
	}
	vm, err := presenter.HandleIntent(ctx, i, v.State())
	if err != nil {
		return err
	}
	return v.Display(ctx, vm)
}
