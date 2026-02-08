package app

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/config"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/controller/http"
	tgbotcontroller "github.com/soumirel/wishlister/services/telegram-bot/internal/controller/telegram_bot"
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

	botHandler := tgbotcontroller.NewBotHandler(wishlisterAuthSvc)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.Handle),
	}

	b, err := bot.New(cfg.TelegramBotConfig.Token, opts...)
	if err != nil {
		panic(err)
	}

	go b.Start(ctx)

	select {}
}
