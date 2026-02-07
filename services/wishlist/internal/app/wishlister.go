package app

import (
	"log"
	"os"

	httpController "github.com/soumirel/wishlister/services/wishlist/internal/controller/http"
	grpcController "github.com/soumirel/wishlister/services/wishlist/internal/controller/grpc"
	"github.com/soumirel/wishlister/services/wishlist/internal/config"
	"github.com/soumirel/wishlister/services/wishlist/internal/repository"
	"github.com/soumirel/wishlister/services/wishlist/internal/uof"
	useruc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/user"
	useridentity "github.com/soumirel/wishlister/services/wishlist/internal/usecase/user_identity"
	wishuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wish"
	wishlistuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wishlist"
	wishlistpermuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wishlist_permission"
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config:", err)
	}

	migrationsBytes, err := os.ReadFile(cfg.Paths.Migrations)
	if err != nil {
		log.Fatal("open migrations file failed:", err.Error())
	}
	refreshScriptBytes, err := os.ReadFile(cfg.Paths.RefreshData)
	if err != nil {
		log.Fatal("open refresh data file failed:", err.Error())
	}

	postgresClient := repository.InitPostgresClient(
		cfg.DbConfig(),
		string(migrationsBytes),
		string(refreshScriptBytes),
	)

	uofFactory := uof.NewUnitOfWorkFactory(postgresClient)

	userUc := useruc.NewUserUsecase(uofFactory)
	wishlistUc := wishlistuc.NewWishlistUsecase(uofFactory)
	wishUc := wishuc.NewWishUsecase(uofFactory)
	wishlistPermissionUc := wishlistpermuc.NewWishlistPermissionUsecase(uofFactory)
	userIdentityUc := useridentity.NewUserIdentityUsecase(uofFactory)

	httpController.StartHttpServer(cfg.Server.HTTPAddr, userUc, wishlistUc, wishUc, wishlistPermissionUc)
	grpcController.StartGrpcServer(cfg.Server.GRPCAddr, userIdentityUc)

	select {}
}
