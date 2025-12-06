package domain

import (
	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID   string
	Name string
}

func NewUser() *User {
	id := uuid.Must(uuid.NewV4())
	return &User{
		ID: id.String(),
	}
}
