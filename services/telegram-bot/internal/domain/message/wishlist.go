package message

import (
	"fmt"
	"strings"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/entity"
)

const (
	EmptyWishlistListMsg = "no wishlists"
)

func MakeGetWishlistsMessage(list entity.WishlistListModel) string {
	if len(list) == 0 {
		return EmptyWishlistListMsg
	}
	msgItems := make([]string, 0, len(list))
	for _, item := range list {
		msgItems = append(msgItems, fmt.Sprintf("ID: %v\nName: %v", item.ID, item.Name))

	}
	return strings.Join(msgItems, "\n\n")
}
