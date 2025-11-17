package controller

import (
	"wishlister/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	serverAddr = ":8080"
)

func StartHttpServer(userService *service.UserService,
	wishService *service.WishService) {
	e := gin.New()

	{
		v1Group := e.Group("/v1")

		userGroup := v1Group.Group("/user")
		_ = NewUserHandler(userGroup, userService)

		wishGroup := v1Group.Group("/wish")
		_ = NewWishHandler(wishGroup, wishService)
	}

	e.Run(serverAddr)
}
