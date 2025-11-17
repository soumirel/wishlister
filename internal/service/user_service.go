package service

import (
	"context"
	"wishlister/internal/domain"
)

type userRepository interface {
	GetList(context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id string) (domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type UserService struct {
	userRepo userRepository
}

func NewUserService(userRepo userRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetList(ctx context.Context) ([]domain.User, error) {
	return s.userRepo.GetList(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id string) (domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, req domain.User) (domain.User, error) {
	user := domain.NewUser()
	user.Name = req.Name
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
