package controller

import (
	v1 "wishlister/internal/controller/v1"

	"github.com/gin-gonic/gin"
)

const (
	serverAddr = ":8080"
)

func StartHttpServer(
	v1WishlistService v1.WishlistService,
	v1WishService v1.WishService,
) {
	e := gin.New()

	{
		v1Group := e.Group("/v1")

		v1.BindHandlers(v1Group, v1WishlistService, v1WishService)
	}

	e.Run(serverAddr)
}
