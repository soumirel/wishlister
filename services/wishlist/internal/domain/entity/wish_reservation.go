package entity

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type WishReservation struct {
	ID               string
	ReservedByUserID string
	ReservedAt       time.Time
}

func newWishReservation(userID string, reservedTime time.Time) *WishReservation {
	id := uuid.Must(uuid.NewV4()).String()
	return &WishReservation{
		ID:               id,
		ReservedByUserID: userID,
		ReservedAt:       reservedTime,
	}
}
