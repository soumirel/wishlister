package wishlist_permission

type GrantWishlistPermissionCommand struct {
	RequestorUserID string
	WishlistID      string
	TargetUserID    string
	PermissionLevel string
}

type RevokeWishlistPermissionCommand struct {
	RequestorUserID string
	WishlistID      string
	TargetUserID    string
}
