package app

import (
	"log"
	"os"
	"wishlister/internal/controller"
	"wishlister/internal/repository"
	"wishlister/internal/uof"
	useruc "wishlister/internal/usecase/user"

	wishuc "wishlister/internal/usecase/wish"
	wishlistuc "wishlister/internal/usecase/wishlist"
	wishlistpermuc "wishlister/internal/usecase/wishlist_permission"
)

func Run() {
	migrationsBytes, err := os.ReadFile("./db/init/init.sql")
	if err != nil {
		log.Fatal("open migrations file failed:", err.Error())
	}
	refreshScriptBytes, err := os.ReadFile("./db/refresh_test_data.sql")
	if err != nil {
		log.Fatal("open migrations file failed:", err.Error())
	}
	postgresClient := repository.InitPostgresClient(string(migrationsBytes), string(refreshScriptBytes))

	uofFactory := uof.NewUnitOfWorkFactory(postgresClient)

	userUc := useruc.NewUserUsecase(uofFactory)
	wishlistUc := wishlistuc.NewWishlistUsecase(uofFactory)
	wishUc := wishuc.NewWishUsecase(uofFactory)
	wishlistPermissionUc := wishlistpermuc.NewWishlistPermissionUsecase(uofFactory)

	controller.StartHttpServer(userUc, wishlistUc, wishUc, wishlistPermissionUc)
}
