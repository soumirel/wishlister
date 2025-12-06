package domain

import (
	"github.com/gofrs/uuid/v5"
)

type Wish struct {
	ID         string
	UserID     string
	WishlistID string
	Name       string
	//ReservedByUserID string
}

func newWish(userID, wishlistID, name string) *Wish {
	id := uuid.Must(uuid.NewV4()).String()
	return &Wish{
		ID:         id,
		UserID:     userID,
		WishlistID: wishlistID,
		Name:       name,
	}
}

// type Reservation struct {
// 	UserID string `json:"userId"`
// 	WishID string `json:"wishId"`
// }
