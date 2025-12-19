package repository

import (
	"context"
	"errors"
	"wishlister/internal/domain/entity"
	"wishlister/internal/domain/repository"

	"github.com/jackc/pgx/v5"
)

type wishRepository struct {
	q Querier
}

func newWishRepository(q Querier) repository.WishRepository {
	return &wishRepository{
		q: q,
	}
}

func (r *wishRepository) GetWish(ctx context.Context, wishlistID, wishID string) (*entity.Wish, error) {
	query := `
			SELECT 
				id, wishlist_id, name
			FROM wishes
			WHERE wishlist_id = $1
				AND id = $2`
	var wish entity.Wish
	err := r.q.QueryRow(ctx, query, wishlistID, wishID).Scan(
		&wish.ID, &wish.WishlistID, &wish.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrWishDoesNotExist
		}
		return nil, err
	}
	return &wish, nil
}

func (r *wishRepository) UpdateWish(ctx context.Context, wish *entity.Wish) error {
	query := `
			UPDATE wishes 
			SET name = $3
			WHERE wishlist_id = $1
				AND id = $2`
	_, err := r.q.Exec(ctx, query, wish.WishlistID, wish.ID, wish.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishRepository) CreateWish(ctx context.Context, wish *entity.Wish) error {
	query := `
			INSERT INTO wishes(id, wishlist_id, name)
			VALUES ($1, $2, $3)`
	_, err := r.q.Exec(ctx, query,
		wish.ID, wish.WishlistID, wish.Name,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishRepository) DeleteWish(ctx context.Context, wishlistID, wishID string) error {
	query := `
			DELETE FROM wishes
			WHERE wishlist_id = $1
				AND id = $2`
	_, err := r.q.Exec(ctx, query, wishlistID, wishID)
	if err != nil {
		return err
	}
	return nil
}
