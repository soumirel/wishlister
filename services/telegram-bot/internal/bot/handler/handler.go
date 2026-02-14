package telegrambot

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/auth"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/view"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

type botHandler struct {
	intentDispatcher ui.IntentDispatcher
}

func NewBotHandler(
	intentDispatcher ui.IntentDispatcher,
) *botHandler {
	return &botHandler{
		intentDispatcher: intentDispatcher,
	}
}

func (h *botHandler) HandleListCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := h.intentDispatcher.DispatchViewIntent(ctx, ui.ShowWishlistsIntent{}, view.NewTelegramView(b, update))
	if err != nil {
		log.Printf("presenter error: %v\n", err.Error())
	}
}

func (h *botHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	au, ok := auth.FromCtx(ctx)
	if !ok {
		log.Print("no auth in handler context")
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("%v from userID: %v", update.Message.Text, au.UserID),
	})
}
