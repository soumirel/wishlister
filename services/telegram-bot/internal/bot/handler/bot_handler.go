package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
)

type botHandler struct {
	viewActivityStore model.ViewStore
	viewFactory       model.ViewRegistry
}

func NewBotHandler(
	viewActivityStore model.ViewStore,
	viewFactory model.ViewRegistry,
) *botHandler {
	return &botHandler{
		viewActivityStore: viewActivityStore,
		viewFactory:       viewFactory,
	}
}

func (d *botHandler) HandleUpdate(ctx context.Context, bot *bot.Bot, update *models.Update) error {
	activeViewType, err := d.viewActivityStore.GetFocusedViewType(ctx, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	var v model.UpdateProcessor
	switch activeViewType {
	case model.HubViewType:
		v, err = d.viewFactory.HubView(ctx)
	case model.WishlistCreatorViewType:
		v, err = d.viewFactory.WishlistCreatorView(ctx)
	case model.WishlistsViewType:
		v, err = d.viewFactory.WishlistsView(ctx)
	}
	return v.ProcessUpdate(ctx, update)
}

func (d *botHandler) HandleCommand(ctx context.Context, bot *bot.Bot, update *models.Update, cmd model.Command) error {
	v, err := d.viewFactory.HubView(ctx)
	if err != nil {
		return err
	}
	err = v.ProcessCommand(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}
