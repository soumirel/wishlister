package repository

import (
	"context"

	"github.com/soumirel/wishlister/telegram-bot/internal/domain/entity"
)

type ExternalIdentity struct {
	ExternalID       string
	IdentityProvider string
}

type WishlistCoreAuthRepository interface {
	GetUserIdByExternalIdentity(ctx context.Context, ei ExternalIdentity) (string, error)
	CreateUserFromExternalIdentity(ctx context.Context, ei ExternalIdentity) (string, error)
}

type WishlisterReadRepository interface {
	GetWishlists(ctx context.Context) ([]*entity.Wishlist, error)
}
