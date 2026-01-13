package service

import (
	"context"

	"github.com/soumirel/wishlister/wishlist/internal/domain/entity"
)

type CreateUserOptions struct {
	Name string
}

type UserService interface {
	GetUsers(ctx context.Context) ([]*entity.User, error)
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	CreateUser(ctx context.Context, o CreateUserOptions) (*entity.User, error)
	DeleteUser(ctx context.Context, userID string) error
}
