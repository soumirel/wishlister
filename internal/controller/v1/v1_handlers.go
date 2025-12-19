package controller

import (
	"wishlister/internal/controller/v1/middleware"
	useruc "wishlister/internal/usecase/user"
	wishuc "wishlister/internal/usecase/wish"
	wishlistuc "wishlister/internal/usecase/wishlist"
	wishlistpermuc "wishlister/internal/usecase/wishlist_permission"

	"github.com/gin-gonic/gin"
)

func BindHandlers(
	gr *gin.RouterGroup,
	userUc *useruc.UserUsecase,
	wishlistUc *wishlistuc.WishlistUsecase,
	wishUc *wishuc.WishUsecase,
	wishlistPermissionUc *wishlistpermuc.WishlistPermissionUsecase,
) {

	wishlistsGr := gr.Group("/wishlists",
		middleware.AuthMiddleware(),
		middleware.ErrorHandler(),
	)
	NewWishlistHandler(wishlistsGr, wishlistUc, wishUc, wishlistPermissionUc)

	usersGr := gr.Group("/users",
		middleware.ErrorHandler(),
	)
	NewUserHandler(usersGr, userUc)
}
