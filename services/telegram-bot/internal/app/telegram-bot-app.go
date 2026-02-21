package app

import (
	"context"
	"log"
	"os"
	"os/signal"

	tgbothandler "github.com/soumirel/wishlister/services/telegram-bot/internal/bot/handler"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/ui"

	tgbotstore "github.com/soumirel/wishlister/services/telegram-bot/internal/bot/store"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/config"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/controller/http"
	grpcrepo "github.com/soumirel/wishlister/services/telegram-bot/internal/repository/grpc"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/repository/store"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/service"
	"github.com/valkey-io/valkey-go"
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config:", err)
	}

	// healthcheck server
	http.StartHttpServer(cfg.Server.HTTPAddr)

	// storages
	valkeyClient, err := valkey.NewClient(
		valkey.ClientOption{
			InitAddress: []string{cfg.Valkey.Addr},
			Password:    cfg.Valkey.Password,
		},
	)
	if err != nil {
		panic(err)
	}

	// grpc clients
	wishlisterGRPC, err := grpcrepo.NewWishlistGrpcRepository(cfg.Services.WishlistGrpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	// services
	wishlisterAuthSvc := service.NewWishlisterAuthSvc(wishlisterGRPC)
	wishlistSvc := service.NewWishlisterSvc(wishlisterGRPC)

	// store clients
	stateStore := store.NewStateStore(valkeyClient)
	viewActivityStore := tgbotstore.NewViewStore(valkeyClient)
	vlsStore := store.NewVlsStore(valkeyClient)

	// store clients wrappers
	wishlistCreatorStateStore := store.NewWishlistCreatorStateStore(stateStore)

	// view lifecycle
	vlc := ui.NewViewLifecycleController(vlsStore)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = tgbothandler.StartTelegramBot(
		ctx,
		cfg.TelegramBot.Token,
		wishlisterAuthSvc,
		viewActivityStore,
		wishlistSvc,
		wishlistCreatorStateStore,
		vlc,
	)
	if err != nil {
		panic(err)
	}

	select {}
}
