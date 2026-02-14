package service

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/repository"
)

type wishlisterReadSvc struct {
	readRepo repository.WishlistCoreReadRepository
}

func NewWishlisterReadSvc(readRepo repository.WishlistCoreReadRepository) *wishlisterReadSvc {
	return &wishlisterReadSvc{
		readRepo: readRepo,
	}
}

func (s *wishlisterReadSvc) GetWishlists(ctx context.Context) (model.WishlistList, error) {
	return s.readRepo.GetWishlists(ctx)
}
