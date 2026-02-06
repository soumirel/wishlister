package telegrambot

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/entity"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
)

const (
	identityProvider = "telegram"
)

type botHandler struct {
	authSvc service.WishlisterAuthService
}

func NewBotHandler(authSvc service.WishlisterAuthService) *botHandler {
	return &botHandler{
		authSvc: authSvc,
	}
}

func (h *botHandler) HandleListCommand(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func (h *botHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	user, err := h.authUpdate(ctx, update)
	if err != nil {
		log.Print(fmt.Errorf("auth error: %w", err))
		return
	}
	if user == nil {
		log.Print("no wishlister user, abort handling")
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("%v from userID: %v", update.Message.Text, user.ID),
	})
}

func (h *botHandler) authUpdate(ctx context.Context, update *models.Update) (*entity.WishlisterUser, error) {
	senderTelegramUserID, ok := getTelegramSenderUserID(update)
	if !ok {
		return nil, errors.New("no sender id in update")
	}
	wishlisterUser, err := h.authSvc.AuthByTelegramID(ctx, senderTelegramUserID)
	if err != nil {
		return nil, err
	}
	return wishlisterUser, nil
}

func getTelegramSenderUserID(update *models.Update) (int64, bool) {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return 0, false
	}
	return update.Message.From.ID, true
}
