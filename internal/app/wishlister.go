package app

import (
	"log"
	"os"
	"wishlister/internal/controller"
	"wishlister/internal/repository"
	wishlist "wishlister/internal/service/wishlist"
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
	dbClient := repository.InitDbClient(string(migrationsBytes), string(refreshScriptBytes))

	wishlistRepo := repository.NewWishlistRepository(dbClient)
	wishRepo := repository.NewWishRepository(dbClient)
	wishlistPermissionRepo := repository.NewWishlistPersmissionRepository(dbClient)

	wishlistService := wishlist.NewWishlistService(
		wishlistRepo, wishRepo, wishlistPermissionRepo,
		dbClient,
	)

	controller.StartHttpServer(wishlistService, wishlistService)
}
