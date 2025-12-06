package dto

type GetWishlistsRequest struct {
}

type GetWishlistsResponse struct {
}

type GetWishlistRequest struct {
}

type CreateWishlistRequest struct {
	WishlistName string `json:"name" binding:"required"`
}

type UpdateWishlistRequest struct {
	WishlistName string `json:"name" binding:"required"`
}

type DeleteWishlistRequest struct {
}

type GrantWishlistPermissionRequest struct {
	UserID          string `json:"user_id" binding:"required"`
	PermissionLevel string `json:"permission_level" binding:"required"`
}

type RevokeWishlistPermissionRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type CreateWishRequest struct {
	WishName string `json:"name" binding:"required"`
}

type UpdateWishRequest struct {
	WishName string `json:"name" binding:"required"`
}
