package domain

import "github.com/gofrs/uuid/v5"

type Wish struct {
	ID     string `json:"id"`
	UserID string `json:"userId`
	Name   string `json: "name"`
}

func NewWish(UserID string) *Wish {
	id := uuid.Must(uuid.NewV4())
	return &Wish{
		ID:     id.String(),
		UserID: UserID,
	}
}
