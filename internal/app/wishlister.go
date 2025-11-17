package app

import (
	"wishlister/internal/controller"
	"wishlister/internal/repository"
	"wishlister/internal/service"
)

func Run() {
	dbPool := repository.InitDB()
	userRepository := repository.NewUserRepository(dbPool)
	wishRepository := repository.NewWishRepository(dbPool)

	userService := service.NewUserService(userRepository)
	wishService := service.NewWishService(wishRepository)

	controller.StartHttpServer(userService, wishService)
}
