package model

import "errors"

var (
	ErrUserDoesNotExist = errors.New("user does not exist")
)

type WishlisterUser struct {
	ID string
}
