package app

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	tgbotcontroller "github.com/soumirel/wishlister/services/telegram-bot/internal/controller/telegram_bot"
	grpcrepo "github.com/soumirel/wishlister/services/telegram-bot/internal/repository/grpc"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/service"
)

const (
	telegramBotTokenEnv = "TELEGRAM_BOT_TOKEN"

	wishlistGrpcAddr = ":8081"
)

func Run() {
	wishlisterGRPC, err := grpcrepo.NewWishlistGrpcRepository(wishlistGrpcAddr)
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

	b, err := bot.New(os.Getenv(telegramBotTokenEnv), opts...)
	if err != nil {
		panic(err)
	}

	go b.Start(ctx)

	log.Print("telegram bot successfuly started")

	select {}
}
