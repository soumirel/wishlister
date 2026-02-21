package view

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	botwrap "github.com/soumirel/go-telegram-null-safety/wrap"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/bot/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"
	wishlistcreator "github.com/soumirel/wishlister/services/telegram-bot/internal/ui/wishlist_creator"
)

const (
	namePromt = "Enter wishlist name"
)

type wishlistCreatorView struct {
	bot    *bot.Bot
	chatID int64

	presenter wishlistcreator.Presenter
}

func NewWishlistCreatorView(
	bot *bot.Bot,
	chatID int64,
	presenter wishlistcreator.Presenter,
) *wishlistCreatorView {
	return &wishlistCreatorView{
		bot:       bot,
		chatID:    chatID,
		presenter: presenter,
	}
}

func (v *wishlistCreatorView) ViewIdentifier() ui.ViewIdentifier {
	return ui.ViewIdentifier(
		fmt.Sprintf("%v:%v", model.WishlistCreatorIdPrefix, v.chatID),
	)
}

func (v *wishlistCreatorView) OnCreate(ctx context.Context) error {
	v.presenter.AttachView(v)
	return nil
}

func (v *wishlistCreatorView) OnStart(ctx context.Context) error {
	return v.presenter.OnStart(ctx)
}

func (v *wishlistCreatorView) OnResume(context.Context) error {
	v.presenter.AttachView(v)
	return nil
}

func (v *wishlistCreatorView) ProcessUpdate(ctx context.Context, update *models.Update) error {
	text := botwrap.Update(update).GetMessage().GetText()
	if text == "" {
		return nil
	}
	err := v.presenter.OnNameEntered(ctx, wishlistcreator.NameEnteringIntent{
		Name: text,
	})
	if err != nil {
		return err
	}
	return nil
}

func (v *wishlistCreatorView) buildNamePromtMessage(vm wishlistcreator.NamePromtVM) string {
	return namePromt
}

func (v *wishlistCreatorView) DisplayNamePromt(ctx context.Context, vm wishlistcreator.NamePromtVM) error {
	_, err := v.bot.SendMessage(ctx,
		&bot.SendMessageParams{
			ChatID: v.chatID,
			Text:   v.buildNamePromtMessage(vm),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (v *wishlistCreatorView) buildWishlistCreatedMessage(vm wishlistcreator.CreatedVM) string {
	return fmt.Sprintf(
		`Wishlist created

ID: %v
Name: %v`,
		vm.WishlistVM.ID, vm.WishlistVM.Name)
}

func (v *wishlistCreatorView) DisplayWishlistCreated(ctx context.Context, vm wishlistcreator.CreatedVM) error {
	_, err := v.bot.SendMessage(ctx,
		&bot.SendMessageParams{
			ChatID: v.chatID,
			Text:   v.buildWishlistCreatedMessage(vm),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
