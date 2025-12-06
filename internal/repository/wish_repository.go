package repository

import (
	"context"
	"errors"
	"wishlister/internal/domain"
	"wishlister/internal/pkg"

	"github.com/jackc/pgx/v5"
)

type wishRepository struct {
	db pkg.DbExecutor
}

func NewWishRepository(db pkg.DbExecutor) *wishRepository {
	return &wishRepository{
		db: db,
	}
}

func (r *wishRepository) GetWish(ctx context.Context, userID, wishlistID, wishID string) (*domain.Wish, error) {
	query := `
			SELECT 
				id, user_id, wishlist_id, name
			FROM wishes
			WHERE user_id = $1
				AND wishlist_id = $2
				AND id = $3`
	var wish domain.Wish
	err := r.db.QueryRow(ctx, query, userID, wishlistID, wishID).Scan(
		&wish.ID, &wish.UserID, &wish.WishlistID, &wish.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrWishDoesNotExist
		}
		return nil, err
	}
	return &wish, nil
}

func (r *wishRepository) UpdateWish(ctx context.Context, wish *domain.Wish) error {
	query := `
			UPDATE wishes 
			SET name = $4
			WHERE user_id = $1
				AND wishlist_id = $2
				AND id = $3`
	_, err := r.db.Exec(ctx, query, wish.UserID, wish.WishlistID, wish.ID, wish.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishRepository) CreateWish(ctx context.Context, wish *domain.Wish) error {
	query := `
			INSERT INTO wishes(id, user_id, wishlist_id, name)
			VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query,
		wish.ID, wish.UserID, wish.WishlistID, wish.Name,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishRepository) DeleteWish(ctx context.Context, userID, wishlistID, wishID string) error {
	query := `
			DELETE FROM wishes
			WHERE user_id = $1
				AND wishlist_id = $2
				AND id = $3`
	_, err := r.db.Exec(ctx, query, userID, wishlistID, wishID)
	if err != nil {
		return err
	}
	return nil
}
