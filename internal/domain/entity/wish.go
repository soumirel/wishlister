package entity

import (
	"github.com/gofrs/uuid/v5"
)

type Wish struct {
	ID         string
	WishlistID string
	Name       string
}

func newWish(wishlistID, name string) *Wish {
	id := uuid.Must(uuid.NewV4()).String()
	return &Wish{
		ID:         id,
		WishlistID: wishlistID,
		Name:       name,
	}
}

func (w *Wish) UpdateName(name string) error {
	w.Name = name
	return nil
}
