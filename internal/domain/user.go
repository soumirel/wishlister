package domain

import (
	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewUser() *User {
	id := uuid.Must(uuid.NewV4())
	return &User{
		ID: id.String(),
	}
}
