package entity

type WishlistListItemModel struct {
	ID     string
	UserID string
	Name   string
}

type WishlistListModel []*WishlistListItemModel
