package repository

import (
	"context"
	"wishlister/internal/domain"

	"github.com/georgysavva/scany/v2/pgxscan"

	"github.com/jackc/pgx/v5/pgxpool"
)

type wishRepository struct {
	db *pgxpool.Pool
}

func NewWishRepository(pool *pgxpool.Pool) *wishRepository {
	return &wishRepository{
		db: pool,
	}
}

func (r *wishRepository) GetList(ctx context.Context) ([]domain.Wish, error) {
	query := `SELECT id, user_id, name FROM wishes`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var wishes []domain.Wish
	err = pgxscan.ScanAll(&wishes, rows)
	if err != nil {
		return nil, err
	}
	return wishes, nil
}

func (r *wishRepository) GetByID(ctx context.Context, id string) (domain.Wish, error) {
	query := `SELECT id, user_id, name FROM wishes WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var wish domain.Wish
	err := row.Scan(&wish.ID, &wish.UserID, &wish.Name)
	if err != nil {
		return domain.Wish{}, err
	}
	return wish, nil
}

func (r *wishRepository) Create(ctx context.Context, wish *domain.Wish) error {
	query := `INSERT INTO wishes(id, user_id, name) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, wish.ID, wish.UserID, wish.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM wishes WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
