package model

import "errors"

type Command string

const (
	WishlistsCommand      Command = "wishlists"
	CreateWishlistCommand Command = "createwishlist"
)

var (
	ErrUnknownCommand = errors.New("unknown command")
)
