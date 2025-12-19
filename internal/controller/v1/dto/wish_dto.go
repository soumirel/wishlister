package dto

type GetWishRequest struct {
}

type CreateWishRequest struct {
	WishName string `json:"name" binding:"required"`
}

type UpdateWishRequest struct {
	WishName string `json:"name" binding:"required"`
}

type DeleteWishRequest struct {
}
