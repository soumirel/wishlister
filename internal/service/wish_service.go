package service

import (
	"context"
	"wishlister/internal/domain"
)

type wishRepository interface {
	GetList(context.Context) ([]domain.Wish, error)
	GetByID(ctx context.Context, id string) (domain.Wish, error)
	Create(ctx context.Context, wish *domain.Wish) error
	Delete(ctx context.Context, id string) error
}

type WishService struct {
	wishRepo wishRepository
}

func NewWishService(wishRepo wishRepository) *WishService {
	return &WishService{
		wishRepo: wishRepo,
	}
}

func (s *WishService) GetList(ctx context.Context) ([]domain.Wish, error) {
	return s.wishRepo.GetList(ctx)
}

func (s *WishService) GetByID(ctx context.Context, id string) (domain.Wish, error) {
	return s.wishRepo.GetByID(ctx, id)
}

func (s *WishService) Create(ctx context.Context, req domain.Wish) (domain.Wish, error) {
	wish := domain.NewWish(req.UserID)
	wish.Name = req.Name
	err := s.wishRepo.Create(ctx, wish)
	if err != nil {
		return domain.Wish{}, err
	}
	return *wish, nil
}

func (s *WishService) Delete(ctx context.Context, id string) error {
	return s.wishRepo.Delete(ctx, id)
}
