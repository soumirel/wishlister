package view

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui/hub"
)

type hubView struct {
	bot    *bot.Bot
	chatID int64

	presenter hub.Presenter
}

func NewHubView(
	bot *bot.Bot,
	chatID int64,
	presenter hub.Presenter,
) *hubView {
	return &hubView{
		bot:       bot,
		chatID:    chatID,
		presenter: presenter,
	}
}

func (v *hubView) ViewIdentifier() ui.ViewIdentifier {
	return ui.ViewIdentifier(
		fmt.Sprintf("%v:%v", model.HubIdPrefix, v.chatID),
	)
}

func (v *hubView) OnCreate(ctx context.Context) error {
	v.presenter.AttachView(v)
	return nil
}

func (v *hubView) OnStart(ctx context.Context) error {
	return nil
}

func (v *hubView) OnResume(context.Context) error {
	v.presenter.AttachView(v)
	return nil
}

func (v *hubView) ProcessUpdate(ctx context.Context, update *models.Update) error {
	return nil
}

func (v *hubView) ProcessCommand(ctx context.Context, cmd model.Command) error {
	switch cmd {
	case model.CreateWishlistCommand:
		return v.presenter.OnNavigateToWishlistCreator(ctx, hub.NavigateToWishlistCreatorIntent{})
	case model.WishlistsCommand:
		return v.presenter.OnNavigateToWishlists(ctx, hub.NavigateToWishlistsIntent{})
	}
	return model.ErrUnknownCommand
}
