package entity

import (
	"errors"

	"github.com/gofrs/uuid/v5"
)

type Wishlist struct {
	ID     string
	UserID string
	Name   string
}

var (
	ErrWishlistDoesNotExist = errors.New("wishlist does not exist")
	ErrWishDoesNotExist     = errors.New("wish does not exist")
)

func NewWishlist(userID, name string) *Wishlist {
	id := uuid.Must(uuid.NewV4()).String()
	return &Wishlist{
		ID:     id,
		UserID: userID,
		Name:   name,
	}
}

func (l *Wishlist) NewWish(wishName string) (*Wish, error) {
	wish := newWish(l.ID, wishName)
	return wish, nil
}
