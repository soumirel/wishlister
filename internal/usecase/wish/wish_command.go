package wish

type GetWishCommand struct {
	RequestorUserID string
	WishlistID      string
	WishID          string
}

type CreateWishCommand struct {
	RequestorUserID string
	WishlistID      string
	WishName        string
}

type UpdateWishCommand struct {
	RequestorUserID string
	WishlistID      string
	WishID          string
	WishName        string
}

type DeleteWishCommand struct {
	RequestorUserID string
	WishlistID      string
	WishID          string
}
