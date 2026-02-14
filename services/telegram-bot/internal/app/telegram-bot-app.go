package app

import (
	"context"
	"log"
	"os"
	"os/signal"

	tgbothandler "github.com/soumirel/wishlister/services/telegram-bot/internal/bot/handler"
	tgbotrepo "github.com/soumirel/wishlister/services/telegram-bot/internal/bot/repository"
	tgbotsvc "github.com/soumirel/wishlister/services/telegram-bot/internal/bot/service"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/config"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/controller/http"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/presenter"
	grpcrepo "github.com/soumirel/wishlister/services/telegram-bot/internal/repository/grpc"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/service"
	"github.com/valkey-io/valkey-go"
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config:", err)
	}

	http.StartHttpServer(cfg.Server.HTTPAddr)

	valkeyClient, err := valkey.NewClient(
		valkey.ClientOption{
			InitAddress: []string{cfg.Valkey.Addr},
			Password:    cfg.Valkey.Password,
		},
	)
	if err != nil {
		panic(err)
	}

	wishlisterGRPC, err := grpcrepo.NewWishlistGrpcRepository(cfg.Services.WishlistGrpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	wishlisterAuthSvc := service.NewWishlisterAuthSvc(wishlisterGRPC)
	wishlistReadSvc := service.NewWishlisterReadSvc(wishlisterGRPC)

	svcFactory := service.NewServiceFactory(
		wishlistReadSvc,
	)

	presenterProvider, err := presenter.NewPresenterProvider(
		presenter.NewMasterPresenter(svcFactory),
	)
	if err != nil {
		panic(err)
	}

	intentDispatcher := presenter.NewIntentDispatcher(presenterProvider)
	viewStateRepo := tgbotrepo.NewViewStateRepository(valkeyClient)
	viewStateSvc := tgbotsvc.NewViewStateService(viewStateRepo)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = tgbothandler.StartTelegramBot(
		ctx, cfg.TelegramBot.Token,
		wishlisterAuthSvc, intentDispatcher,
		viewStateSvc,
	)
	if err != nil {
		panic(err)
	}

	select {}
}
