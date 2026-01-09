package controller

import (
	"github.com/soumirel/wishlister/wishlist/internal/controller/v1/middleware"
	useruc "github.com/soumirel/wishlister/wishlist/internal/usecase/user"
	wishuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wish"
	wishlistuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wishlist"
	wishlistpermuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wishlist_permission"

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
