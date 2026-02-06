package repository

import (
	"context"

	"github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
)

type UserIdentityRepository interface {
	GetUserIdByExternalIdentity(ctx context.Context, externalIdentity entity.ExternalIdentity) (string, error)
	SaveIdentity(ctx context.Context, userIDentity *entity.UserIdentity) error
}
