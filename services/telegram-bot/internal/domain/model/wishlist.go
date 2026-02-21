package model

type Wishlist struct {
	ID     string
	UserID string
	Name   string
}

type WishlistList []*Wishlist
