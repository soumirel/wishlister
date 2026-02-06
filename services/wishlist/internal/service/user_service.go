package service

import (
	"context"

	"github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/services/wishlist/internal/domain/repository"
	"github.com/soumirel/wishlister/services/wishlist/internal/domain/service"
)

type userService struct {
	userRepo repository.UserRepository
}

func newUserService(userRepo repository.UserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUsers(ctx context.Context) ([]*entity.User, error) {
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, entity.ErrUserDoesNotExist
	}
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, o service.CreateUserOptions) (*entity.User, error) {
	user := entity.NewUser()
	user.Name = o.Name

	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
