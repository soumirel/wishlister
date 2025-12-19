package controller

import (
	v1 "wishlister/internal/controller/v1"

	useruc "wishlister/internal/usecase/user"
	wishuc "wishlister/internal/usecase/wish"
	wishlistuc "wishlister/internal/usecase/wishlist"
	wishlistpermuc "wishlister/internal/usecase/wishlist_permission"

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
