package controller

import (
	v1 "github.com/soumirel/wishlister/wishlist/internal/controller/v1"

	useruc "github.com/soumirel/wishlister/wishlist/internal/usecase/user"
	wishuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wish"
	wishlistuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wishlist"
	wishlistpermuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wishlist_permission"

	"github.com/gin-gonic/gin"
)

const (
	serverAddr = ":8080"
)

func StartHttpServer(
	userUc *useruc.UserUsecase,
	wishlistUc *wishlistuc.WishlistUsecase,
	wishUc *wishuc.WishUsecase,
	wishlistPermissionUc *wishlistpermuc.WishlistPermissionUsecase,
) {
	e := gin.New()

	{
		v1Group := e.Group("/v1")

		v1.BindHandlers(v1Group, userUc, wishlistUc, wishUc, wishlistPermissionUc)
	}

	e.Run(serverAddr)
}
