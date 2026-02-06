package service

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/entity"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/repository"
)

type wishlisterReadSvc struct {
	readRepo repository.WishlisterReadRepository
}

func NewWishlisterReadSvc(readRepo repository.WishlisterReadRepository) *wishlisterReadSvc {
	return &wishlisterReadSvc{
		readRepo: readRepo,
	}
}

func (s *wishlisterReadSvc) GetWishlists(ctx context.Context) ([]*entity.Wishlist, error) {
	return s.readRepo.GetWishlists(ctx)
}
