package service

import (
	"context"

	"github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/wishlist/internal/domain/repository"
)

type userIdentityService struct {
	userIdentityRepo repository.UserIdentityRepository
}

func newUserIdentityServivce(userIdentityRepo repository.UserIdentityRepository) *userIdentityService {
	return &userIdentityService{
		userIdentityRepo: userIdentityRepo,
	}
}

func (s *userIdentityService) GetUserIdByExternalIdentity(ctx context.Context, externalIdentity entity.ExternalIdentity) (string, error) {
	userID, err := s.userIdentityRepo.GetUserIdByExternalIdentity(ctx, externalIdentity)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *userIdentityService) LinkUserWithExternalIdentity(ctx context.Context, userID string, externalIdentity entity.ExternalIdentity) error {
	identity := entity.NewUserIdentity(userID, externalIdentity)
	err := s.userIdentityRepo.SaveIdentity(ctx, identity)
	if err != nil {
		return err
	}
	return nil
}
