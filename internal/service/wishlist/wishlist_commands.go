package wishlist

type GetWishlistsCommand struct {
	RequestorUserID string
}

type GetWishlistCommand struct {
	RequestorUserID string
	WishlistID      string
}

type CreateWishlistCommand struct {
	RequestorUserID string
	Name            string
}

type UpdateWishlistCommand struct {
	RequestorUserID string
	WishlistID      string
	Name            string
}

type DeleteWishlistCommand struct {
	RequestorUserID string
	WishlistID      string
}

type GrantWishlistPermissionCommand struct {
	RequestorUserID  string
	WishlistID       string
	RequestingUserID string
	PersmissionLevel string
}

type RevokeWishlistPermissionCommand struct {
	RequestorUserID  string
	WishlistID       string
	RequestingUserID string
}

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
