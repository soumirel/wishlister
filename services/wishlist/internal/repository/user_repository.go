package repository

import (
	"context"
	"errors"

	entity "github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/services/wishlist/internal/domain/repository"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	q Querier
}

func newUserRepository(q Querier) repository.UserRepository {
	return &userRepository{
		q: q,
	}
}

func (r *userRepository) GetUsers(ctx context.Context) ([]*entity.User, error) {
	query := `SELECT id, name FROM users`
	rows, err := r.q.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var users []*entity.User
	err = pgxscan.ScanAll(&users, rows)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUser(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, name FROM users WHERE id = $1`
	row := r.q.QueryRow(ctx, query, id)
	var user entity.User
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserDoesNotExist
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users(id, name) VALUES ($1, $2)`
	_, err := r.q.Exec(ctx, query, user.ID, user.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.q.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
