package controller

import (
	"wishlister/internal/controller/v1/middleware"

	"github.com/gin-gonic/gin"
)

func BindHandlers(
	gr *gin.RouterGroup,
	wishlistService WishlistService,
	wishSerivce WishService) {

	wishlistGr := gr.Group("/wishlists",
		middleware.AuthMiddleware(),
		middleware.ErrorHandler(),
	)
	NewWishlistHandler(wishlistGr, wishlistService, wishSerivce)
}
