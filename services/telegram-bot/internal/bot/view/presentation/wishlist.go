package presentation

import (
	"fmt"
	"strings"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

const (
	EmptyWishlistListMsg = "no wishlists"
)

func MakeShowWishlistsMessage(vm ui.ShowWishlistsViewModel) string {
	if len(vm.Wishlists) == 0 {
		return EmptyWishlistListMsg
	}
	msgItems := make([]string, 0, len(vm.Wishlists))
	for _, item := range vm.Wishlists {
		msgItems = append(msgItems, fmt.Sprintf("ID: %v\nName: %v", item.ID, item.Name))

	}
	return strings.Join(msgItems, "\n\n")
}
