package repository

import (
	"context"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
)

type WishlistCoreAuthRepository interface {
	GetUserIdByExternalIdentity(ctx context.Context, ei model.ExternalIdentity) (string, error)
	CreateUserFromExternalIdentity(ctx context.Context, ei model.ExternalIdentity) (string, error)
}

type WishlistCoreReadRepository interface {
	GetWishlists(ctx context.Context) (model.WishlistList, error)
}
