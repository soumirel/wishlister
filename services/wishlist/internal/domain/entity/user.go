package entity

import (
	"errors"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrUserDoesNotExist = errors.New("user does not exist")
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
