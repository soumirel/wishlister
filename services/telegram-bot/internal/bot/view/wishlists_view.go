package view

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	ui_wishlists "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlists"
)

const (
	EmptyWishlistListMsg = "no wishlists"
)

type wishlistsView struct {
	bot       *bot.Bot
	chatID    int64
	presenter ui_wishlists.Presenter
}

func NewWishlistsView(
	bot *bot.Bot,
	chatID int64,
	presenter ui_wishlists.Presenter,
) *wishlistsView {
	return &wishlistsView{
		bot:       bot,
		chatID:    chatID,
		presenter: presenter,
	}
}

func (v *wishlistsView) ViewIdentifier() ui.ViewIdentifier {
	return ui.ViewIdentifier(
		fmt.Sprintf("%v:%v", model.WishlistsIdPrefix, v.chatID),
	)
}

func (v *wishlistsView) OnCreate(ctx context.Context) error {
	v.presenter.AttachView(v)
	return nil
}

func (v *wishlistsView) OnStart(ctx context.Context) error {
	return v.presenter.OnStart(ctx)
}

func (v *wishlistsView) OnResume(context.Context) error {
	v.presenter.AttachView(v)
	return nil
}

func (v *wishlistsView) ProcessUpdate(ctx context.Context, update *models.Update) error {
	return nil
}

func (v *wishlistsView) buildWishlistsMessage(vm ui_wishlists.WishlistsVM) string {
	if len(vm.Wishlists) == 0 {
		return EmptyWishlistListMsg
	}
	msgItems := make([]string, 0, len(vm.Wishlists))
	for _, item := range vm.Wishlists {
		msgItems = append(msgItems, fmt.Sprintf("ID: %v\nName: %v", item.ID, item.Name))

	}
	return strings.Join(msgItems, "\n\n")
}

func (v *wishlistsView) DisplayWishlists(ctx context.Context, vm ui_wishlists.WishlistsVM) error {
	_, err := v.bot.SendMessage(ctx,
		&bot.SendMessageParams{
			ChatID: v.chatID,
			Text:   v.buildWishlistsMessage(vm),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
