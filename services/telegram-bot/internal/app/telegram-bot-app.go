package app

import (
	"context"
	"log"
	"os"
	"os/signal"

	tgbotcontroller "github.com/soumirel/wishlister/services/telegram-bot/internal/bot/handler"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/config"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/controller/http"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/presenter"
	grpcrepo "github.com/soumirel/wishlister/services/telegram-bot/internal/repository/grpc"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/service"
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config:", err)
	}

	http.StartHttpServer(cfg.Server.HTTPAddr)

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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = tgbotcontroller.StartTelegramBot(
		ctx, cfg.TelegramBot.Token,
		wishlisterAuthSvc, intentDispatcher,
	)
	if err != nil {
		panic(err)
	}

	select {}
}
