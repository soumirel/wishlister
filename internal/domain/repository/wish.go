package repository

import (
	"context"
	entity "wishlister/internal/domain/entity"
)

type WishRepository interface {
	GetWish(ctx context.Context, wishlistID, wishID string) (*entity.Wish, error)
	CreateWish(ctx context.Context, wish *entity.Wish) error
	UpdateWish(ctx context.Context, wish *entity.Wish) error
	DeleteWish(ctx context.Context, wishlistID, wishID string) error
}
