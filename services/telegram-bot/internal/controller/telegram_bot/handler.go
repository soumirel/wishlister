package telegrambot

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/auth"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/message"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
)

const (
	identityProvider = "telegram"
)

type botHandler struct {
	wishlistReadSvc service.WishlistCoreReadService
}

func NewBotHandler(
	wishlistReadSvc service.WishlistCoreReadService,
) *botHandler {
	return &botHandler{
		wishlistReadSvc: wishlistReadSvc,
	}
}

func (h *botHandler) HandleListCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	list, err := h.wishlistReadSvc.GetWishlists(ctx)
	if err != nil {
		log.Print(err.Error())
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message.MakeGetWishlistsMessage(list),
	})
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
