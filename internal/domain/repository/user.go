package repository

import (
	"context"
	entity "wishlister/internal/domain/entity"
)

type UserRepository interface {
	GetUsers(context.Context) ([]*entity.User, error)
	GetUser(ctx context.Context, id string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
}
