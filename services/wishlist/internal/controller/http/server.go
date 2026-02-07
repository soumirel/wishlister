package controller

import (
	"net/http"

	v1 "github.com/soumirel/wishlister/services/wishlist/internal/controller/http/v1"

	useruc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/user"
	wishuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wish"
	wishlistuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wishlist"
	wishlistpermuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wishlist_permission"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func StartHttpServer(
	httpAddr string,
	userUc *useruc.UserUsecase,
	wishlistUc *wishlistuc.WishlistUsecase,
	wishUc *wishuc.WishUsecase,
	wishlistPermissionUc *wishlistpermuc.WishlistPermissionUsecase,
) {
	e := gin.New()

	{
		e.GET("/health", func(ctx *gin.Context) {
			ctx.Render(http.StatusOK, render.String{})
		})
	}

	{
		v1Group := e.Group("/v1")

		v1.BindHandlers(v1Group, userUc, wishlistUc, wishUc, wishlistPermissionUc)
	}

	go func() {
		err := e.Run(httpAddr)
		if err != nil {
			panic(err)
		}
	}()
}
