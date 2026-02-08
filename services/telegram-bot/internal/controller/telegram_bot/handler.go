package telegrambot

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/auth"
)

const (
	identityProvider = "telegram"
)

type botHandler struct {
}

func NewBotHandler() *botHandler {
	return &botHandler{}
}

func (h *botHandler) HandleListCommand(ctx context.Context, b *bot.Bot, update *models.Update) {

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
