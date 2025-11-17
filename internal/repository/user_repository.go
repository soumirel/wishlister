package repository

import (
	"context"
	"wishlister/internal/domain"

	"github.com/georgysavva/scany/v2/pgxscan"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{
		db: pool,
	}
}

func (r *userRepository) GetList(ctx context.Context) ([]domain.User, error) {
	query := `SELECT id, name FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var users []domain.User
	err = pgxscan.ScanAll(&users, rows)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	query := `SELECT id, name FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var user domain.User
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users(id, name) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
