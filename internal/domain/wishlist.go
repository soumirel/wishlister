package domain

import (
	"errors"

	"github.com/gofrs/uuid/v5"
)

type Wishlist struct {
	ID     string
	UserID string
	Name   string
	Wishes map[string]*Wish
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

func (l *Wishlist) AddWish(wishName string) (*Wish, error) {
	wish := newWish(l.UserID, l.ID, wishName)
	l.Wishes[wish.ID] = wish
	return wish, nil
}

func (l *Wishlist) UpdateWish(wishID, wishName string) (*Wish, error) {
	wish, ok := l.Wishes[wishID]
	if !ok {
		return nil, ErrWishDoesNotExist
	}
	wish.Name = wishName
	return wish, nil
}

func (l *Wishlist) DeleteWish(wishID string) error {
	_, ok := l.Wishes[wishID]
	if !ok {
		return ErrWishDoesNotExist
	}
	delete(l.Wishes, wishID)
	return nil
}
