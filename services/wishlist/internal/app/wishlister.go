package app

import (
	"log"
	"os"

	httpController "github.com/soumirel/wishlister/services/wishlist/internal/controller/http"

	grpcController "github.com/soumirel/wishlister/services/wishlist/internal/controller/grpc"
	"github.com/soumirel/wishlister/services/wishlist/internal/repository"
	"github.com/soumirel/wishlister/services/wishlist/internal/uof"
	useruc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/user"
	useridentity "github.com/soumirel/wishlister/services/wishlist/internal/usecase/user_identity"

	wishuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wish"
	wishlistuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wishlist"
	wishlistpermuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wishlist_permission"
)

func Run() {
	migrationsBytes, err := os.ReadFile("./init.sql")
	if err != nil {
		log.Fatal("open migrations file failed:", err.Error())
	}
	refreshScriptBytes, err := os.ReadFile("./refresh_test_data.sql")
	if err != nil {
		log.Fatal("open migrations file failed:", err.Error())
	}
	postgresClient := repository.InitPostgresClient(string(migrationsBytes), string(refreshScriptBytes))

	uofFactory := uof.NewUnitOfWorkFactory(postgresClient)

	userUc := useruc.NewUserUsecase(uofFactory)
	wishlistUc := wishlistuc.NewWishlistUsecase(uofFactory)
	wishUc := wishuc.NewWishUsecase(uofFactory)
	wishlistPermissionUc := wishlistpermuc.NewWishlistPermissionUsecase(uofFactory)
	userIdentityUc := useridentity.NewUserIdentityUsecase(uofFactory)

	httpController.StartHttpServer(userUc, wishlistUc, wishUc, wishlistPermissionUc)
	grpcController.StartGrpcServer(userIdentityUc)

	select {}
}
