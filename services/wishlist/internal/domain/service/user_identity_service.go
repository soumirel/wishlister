package service

import (
	"context"

	"github.com/soumirel/wishlister/wishlist/internal/domain/entity"
)

type UserIdentityService interface {
	GetUserIdByExternalIdentity(ctx context.Context, externalIdentity entity.ExternalIdentity) (string, error)
	LinkUserWithExternalIdentity(ctx context.Context, userID string, externalIdentity entity.ExternalIdentity) error
}
