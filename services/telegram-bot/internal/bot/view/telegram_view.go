package view

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/view/presentation"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

type AfterDisplayCallback func(s model.State)

type telegramView struct {
	botApi            *bot.Bot
	state             model.State
	displayedCallback AfterDisplayCallback
}

func NewTelegramView(
	botApi *bot.Bot,
	state model.State,
	displayedCallback AfterDisplayCallback,
) *telegramView {
	return &telegramView{
		botApi:            botApi,
		state:             state,
		displayedCallback: displayedCallback,
	}
}

func (v *telegramView) Module() ui.ModuleType {
	return ui.MasterModule
}

func (v *telegramView) State() ui.State {
	return v.state
}

func (v *telegramView) Display(ctx context.Context, vm ui.ViewModel) error {
	switch tvm := vm.(type) {
	case ui.ShowWishlistsViewModel:
		return v.displayShowVishlistsVM(ctx, tvm)
	case ui.CreateWishlistNameWaitingVM:
		return v.displayCreateWishlistNameWaitingVM(ctx, tvm)
	default:
		return ui.ErrUnknownViewModel
	}
}

func (v *telegramView) displayShowVishlistsVM(ctx context.Context, vm ui.ShowWishlistsViewModel) error {
	v.botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: v.state.ChatID,
		Text:   presentation.MakeShowWishlistsMessage(vm),
	})
	return nil
}

func (v *telegramView) displayCreateWishlistNameWaitingVM(ctx context.Context, vm ui.CreateWishlistNameWaitingVM) error {
	v.botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: v.state.ChatID,
		Text:   "Enter wishlist name",
	})
	v.state.Type = ui.WishlistCreationNameWaiting
	v.displayedCallback(v.state)
	return nil
}
