package user

import (
	"context"

	"github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/wishlist/internal/domain/repository"
	svcentity "github.com/soumirel/wishlister/wishlist/internal/domain/service"
	"github.com/soumirel/wishlister/wishlist/internal/service"
	"github.com/soumirel/wishlister/wishlist/internal/usecase"
)

type UserUsecase struct {
	uofFactory usecase.UnitOfWorkFactory
}

func NewUserUsecase(
	uofFactory usecase.UnitOfWorkFactory,
) *UserUsecase {
	return &UserUsecase{
		uofFactory: uofFactory,
	}
}

func (uc *UserUsecase) GetUsers(ctx context.Context, cmd GetUsersCommand) ([]entity.User, error) {
	uof := uc.uofFactory.NewUnitOfWork(true)
	var usersRes []entity.User
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		userService := service.NewServiceFactory(rf).UserService()
		users, err := userService.GetUsers(ctx)
		if err != nil {
			return err
		}
		usersRes = make([]entity.User, 0, len(users))
		for _, user := range users {
			usersRes = append(usersRes, *user)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return usersRes, nil
}

func (uc *UserUsecase) GetUser(ctx context.Context, cmd GetUserCommand) (entity.User, error) {
	uof := uc.uofFactory.NewUnitOfWork(true)
	var userRes entity.User
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		userRepo := rf.UserRepository()
		user, err := userRepo.GetUser(ctx, cmd.UserID)
		if err != nil {
			return err
		}
		if user == nil {
			return entity.ErrUserDoesNotExist
		}
		userRes = *user
		return nil
	})
	if err != nil {
		return entity.User{}, err
	}
	return userRes, nil
}

func (uc *UserUsecase) CreateUser(ctx context.Context, cmd CreateUserCommand) (entity.User, error) {
	uof := uc.uofFactory.NewUnitOfWork(true)
	var userRes entity.User
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		userService := service.NewServiceFactory(rf).UserService()
		user, err := userService.CreateUser(ctx, svcentity.CreateUserOptions{
			Name: cmd.Name,
		})
		if err != nil {
			return err
		}
		userRes = *user
		return nil
	})
	if err != nil {
		return entity.User{}, err
	}
	return userRes, nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, cmd DeleteUserCommand) error {
	uof := uc.uofFactory.NewUnitOfWork(true)
	err := uof.Do(ctx, func(ctx context.Context, rf repository.RepositoryFactory) error {
		userService := service.NewServiceFactory(rf).UserService()
		err := userService.DeleteUser(ctx, cmd.UserID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return err
}
