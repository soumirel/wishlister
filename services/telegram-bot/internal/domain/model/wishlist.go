package model

type WishlistListItem struct {
	ID     string
	UserID string
	Name   string
}

type WishlistList []*WishlistListItem
