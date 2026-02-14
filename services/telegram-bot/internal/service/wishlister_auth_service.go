package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/repository"
)

const (
	identityProvider = "telegram"
)

type wishlisterAuthSvc struct {
	wishlisterAuthRepo repository.WishlistCoreAuthRepository
}

func NewWishlisterAuthSvc(wishlisterAuthRepo repository.WishlistCoreAuthRepository) *wishlisterAuthSvc {
	return &wishlisterAuthSvc{
		wishlisterAuthRepo: wishlisterAuthRepo,
	}
}

func (s *wishlisterAuthSvc) AuthByTelegramID(ctx context.Context, telegramID int64) (*model.WishlisterUser, error) {
	ei := model.ExternalIdentity{
		ExternalID:       strconv.FormatInt(telegramID, 10),
		IdentityProvider: identityProvider,
	}
	userID, err := s.wishlisterAuthRepo.GetUserIdByExternalIdentity(ctx, ei)
	if err != nil {
		if errors.Is(err, model.ErrUserDoesNotExist) {
			return s.makeUserResult(s.createUserFromExternalIdentity(ctx, ei))
		}
		return nil, err
	}
	return s.makeUserResult(userID, err)
}

func (s *wishlisterAuthSvc) createUserFromExternalIdentity(ctx context.Context, ei model.ExternalIdentity) (string, error) {
	userID, err := s.wishlisterAuthRepo.CreateUserFromExternalIdentity(ctx, ei)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *wishlisterAuthSvc) makeUserResult(userID string, err error) (*model.WishlisterUser, error) {
	if err != nil {
		return nil, err
	}
	return &model.WishlisterUser{
		ID: userID,
	}, nil
}
