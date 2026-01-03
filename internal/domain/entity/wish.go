package entity

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrWishAlredyReserved = errors.New("wish already reserved")
)

type Wish struct {
	ID          string
	WishlistID  string
	Name        string
	Reservation *WishReservation
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

func (w *Wish) CheckCanRerserve() error {
	if w.IsReserved() {
		return ErrWishAlredyReserved
	}
	return nil
}

func (w *Wish) Reserve(userID string, reservedTime time.Time) error {
	if err := w.CheckCanRerserve(); err != nil {
		return err
	}
	reservation := newWishReservation(
		userID,
		reservedTime,
	)
	w.Reservation = reservation
	return nil
}

func (w *Wish) IsReserved() bool {
	return w.Reservation != nil
}
