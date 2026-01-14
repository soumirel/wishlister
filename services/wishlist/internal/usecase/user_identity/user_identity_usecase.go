package useridentity

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/wishlist/internal/domain/repository"
	"github.com/soumirel/wishlister/wishlist/internal/service"
	"github.com/soumirel/wishlister/wishlist/internal/usecase"
)

type UserIdentityUsecase struct {
	uofFactory usecase.UnitOfWorkFactory
}

func NewUserIdentityUsecase(uofFactory usecase.UnitOfWorkFactory) *UserIdentityUsecase {
	return &UserIdentityUsecase{
		uofFactory: uofFactory,
	}
}

func (uc *UserIdentityUsecase) GetUserIdByExternalIdentity(ctx context.Context, cmd GetUserIdByExternalIdentityCommand) (string, error) {
	externalIdentity, err := entity.NewExternalIdentity(
		cmd.ExternalID, entity.IdentityProvider(cmd.IdentityProvider),
	)
	if err != nil {
		return "", err
	}
	uof := uc.uofFactory.NewUnitOfWork(false)
	var (
		userID string
	)
	err = uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {

		svcFactory := service.NewServiceFactory(rf)
		userID, err = svcFactory.UserIdentityService().GetUserIdByExternalIdentity(ctx, externalIdentity)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (uc *UserIdentityUsecase) LinkUserWithExternalIdentity(ctx context.Context, cmd LinkUserWithExternalIdentityCommand) error {
	identity, err := entity.NewExternalIdentity(cmd.ExternalID, entity.IdentityProvider(cmd.IdentityProvider))
	if err != nil {
		return err
	}
	uof := uc.uofFactory.NewUnitOfWork(true)
	err = uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		svcFactory := service.NewServiceFactory(rf)
		// TODO: use CheckUserExist instead
		_, err := svcFactory.UserService().GetUser(ctx, cmd.UserID)
		if err != nil {
			return nil
		}
		err = svcFactory.UserIdentityService().LinkUserWithExternalIdentity(ctx,
			cmd.UserID, identity,
		)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserIdentityUsecase) CreateUserFromExternalIdentity(ctx context.Context, cmd CreateUserFromExternalIdentityCommand) (string, error) {
	uof := uc.uofFactory.NewUnitOfWork(true)
	var (
		userID string
	)
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		user := entity.NewUser()
		err := user.UpdateName(uuid.Must(uuid.NewV4()).String()[:8])
		if err != nil {
			return err
		}
		err = rf.UserRepository().CreateUser(ctx, user)
		if err != nil {
			return err
		}
		externalIdentity, err := entity.NewExternalIdentity(cmd.ExternalID, entity.IdentityProvider(cmd.IdentityProvider))
		if err != nil {
			return err
		}
		userIdentity := entity.NewUserIdentity(user.ID, externalIdentity)
		err = rf.UserIdentityRepository().SaveIdentity(ctx, userIdentity)
		if err != nil {
			return err
		}
		userID = user.ID
		return nil
	})
	if err != nil {
		return "", err
	}
	return userID, nil
}
