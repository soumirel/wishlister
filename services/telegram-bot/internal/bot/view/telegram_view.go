package view

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/presentation"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

type telegramView struct {
	botApi *bot.Bot
	update *models.Update
}

func NewTelegramView(
	botApi *bot.Bot,
	update *models.Update,
) *telegramView {
	return &telegramView{
		botApi: botApi,
		update: update,
	}
}

func (v *telegramView) Module() ui.ModuleType {
	return ui.MasterModule
}

func (v *telegramView) Display(ctx context.Context, vm ui.ViewModel) error {
	switch tvm := vm.(type) {
	case ui.ShowWishlistsViewModel:
		return v.displayShowVishlistsVM(ctx, tvm)
	default:
		return ui.ErrUnknownViewModel
	}
}

func (v *telegramView) displayShowVishlistsVM(ctx context.Context, vm ui.ShowWishlistsViewModel) error {
	v.botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: v.update.Message.Chat.ID,
		Text:   presentation.MakeShowWishlistsMessage(vm),
	})
	return nil
}
