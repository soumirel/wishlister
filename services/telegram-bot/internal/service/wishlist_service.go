package service

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/repository"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/service"
)

type wishlisterReadSvc struct {
	repo repository.WishlistCoreRepository
}

func NewWishlisterSvc(
	repo repository.WishlistCoreRepository,
) *wishlisterReadSvc {
	return &wishlisterReadSvc{
		repo: repo,
	}
}

func (s *wishlisterReadSvc) GetWishlists(ctx context.Context) (model.WishlistList, error) {
	return s.repo.GetWishlists(ctx)
}

func (s *wishlisterReadSvc) CreateWishlist(ctx context.Context, params service.CreateWishlistParams) (model.Wishlist, error) {
	wishlist, err := s.repo.CreateWishlist(ctx, model.Wishlist{
		Name: params.Name,
	})
	if err != nil {
		return model.Wishlist{}, err
	}
	return wishlist, nil
}
