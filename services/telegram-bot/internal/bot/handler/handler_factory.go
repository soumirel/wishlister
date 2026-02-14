package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	tgmodels "github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/auth"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/view"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

type handlerPattern string

const (
	wishlistsCommand      handlerPattern = "wishlists"
	createWishlistCommand handlerPattern = "createwishlist"
)

type viewStateService interface {
	GetState(ctx context.Context, chatID int64) (model.State, error)
	SaveState(ctx context.Context, s model.State) error
}

type handlerFactory struct {
	intentDispatcher ui.IntentDispatcher
	viewStateSvc     viewStateService
}

func NewHandlerFactory(
	intentDispatcher ui.IntentDispatcher,
	viewStateSvc viewStateService,
) *handlerFactory {
	return &handlerFactory{
		intentDispatcher: intentDispatcher,
		viewStateSvc:     viewStateSvc,
	}
}

func (hf *handlerFactory) RegisterHandler(
	b *bot.Bot,
	handlerType bot.HandlerType,
	pattern handlerPattern,
	matchType bot.MatchType,
	m ...bot.Middleware,
) string {
	var (
		handlerFunc bot.HandlerFunc
	)

	if pattern == emptyPattern {
		handlerFunc = hf.NewBotHandlerFunc()
	} else if handlerType == bot.HandlerTypeCallbackQueryData {
		handlerFunc = hf.NewBotCallbackQueryHandlerFunc()
	} else if matchType == bot.MatchTypeCommand ||
		matchType == bot.MatchTypeCommandStartOnly {
		handlerFunc = hf.NewBotCommandHandlerFunc(pattern)
	} else {
		handlerFunc = hf.NewBotHandlerFunc()
	}

	return b.RegisterHandler(
		handlerType,
		string(pattern),
		matchType,
		handlerFunc,
		m...,
	)
}

func (hf *handlerFactory) getEchoHandleFunc() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *tgmodels.Update) {
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
}

func (hf *handlerFactory) defineIntent(update *models.Update) ui.Intent {
	// TODO
	if update == nil {
		return nil
	}
	return nil
}

func (h *handlerFactory) getAfterViewDisplayCallback() view.AfterDisplayCallback {
	return func(s model.State) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		h.viewStateSvc.SaveState(ctx, s)
	}
}

func (hf *handlerFactory) NewBotHandlerFunc() bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgmodels.Update) {
		state, err := hf.viewStateSvc.GetState(ctx, update.Message.Chat.ID)
		if err != nil {
			log.Printf("get view state failed: %v\n", err.Error())
			return
		}
		intent := hf.defineIntent(update)
		if intent == nil {
			log.Printf("cannot define intent\n")
			return
		}
		view := view.NewTelegramView(bot, state, hf.getAfterViewDisplayCallback())
		err = hf.intentDispatcher.DispatchViewIntent(
			ctx, intent, view,
		)
	}
}

func (hf *handlerFactory) NewBotCallbackQueryHandlerFunc() bot.HandlerFunc {
	// TODO
	return func(ctx context.Context, bot *bot.Bot, update *tgmodels.Update) {
		log.Printf("got callback query, not currently supported\n")
		return
	}
}

func (hf *handlerFactory) defineCommandIntent(commandPattern handlerPattern) ui.Intent {
	switch commandPattern {
	case wishlistsCommand:
		return ui.ShowWishlistsIntent{}
	case createWishlistCommand:
		return ui.CreateWishlistIntent{}
	}
	return nil
}

func (hf *handlerFactory) NewBotCommandHandlerFunc(commandPattern handlerPattern) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *tgmodels.Update) {
		state, err := hf.viewStateSvc.GetState(ctx, update.Message.Chat.ID)
		if err != nil {
			log.Printf("get view state failed: %v\n", err.Error())
			return
		}
		intent := hf.defineCommandIntent(commandPattern)
		if intent == nil {
			log.Printf("cannot define intent\n")
			return
		}
		view := view.NewTelegramView(bot, state, hf.getAfterViewDisplayCallback())
		err = hf.intentDispatcher.DispatchViewIntent(
			ctx, intent, view,
		)
	}
}
