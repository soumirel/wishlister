package app

import (
	"log"
	"os"

	"github.com/soumirel/wishlister/pkg/logger"
	"github.com/soumirel/wishlister/services/wishlist/internal/config"
	grpcController "github.com/soumirel/wishlister/services/wishlist/internal/controller/grpc"
	httpController "github.com/soumirel/wishlister/services/wishlist/internal/controller/http"
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
		log.Fatalf("config loading failed: %v", err.Error())
	}

	logger := logger.Init(map[string]any{
		"service": "wishlist",
	}).Sugar()

	migrationsBytes, err := os.ReadFile(cfg.Paths.Migrations)
	if err != nil {
		logger.Fatalw("migrations file reading failed", err)
	}

	postgresClient := repository.InitPostgresClient(
		cfg.DbConfig(),
		string(migrationsBytes),
	)

	uofFactory := uof.NewUnitOfWorkFactory(postgresClient)

	userUc := useruc.NewUserUsecase(uofFactory)
	wishlistUc := wishlistuc.NewWishlistUsecase(uofFactory)
	wishUc := wishuc.NewWishUsecase(uofFactory)
	wishlistPermissionUc := wishlistpermuc.NewWishlistPermissionUsecase(uofFactory)
	userIdentityUc := useridentity.NewUserIdentityUsecase(uofFactory)

	httpController.StartHttpServer(cfg.Server.HTTPAddr, userUc, wishlistUc, wishUc, wishlistPermissionUc)
	grpcController.StartGrpcServer(cfg.Server.GRPCAddr, userIdentityUc, wishlistUc)

	select {}
}
