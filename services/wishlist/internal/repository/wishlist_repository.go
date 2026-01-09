package repository

import (
	"context"

	"github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/wishlist/internal/domain/repository"

	"github.com/georgysavva/scany/v2/pgxscan"
)

type wishlistRepository struct {
	q Querier
}

func newWishlistRepository(q Querier) repository.WishlistRepository {
	return &wishlistRepository{
		q: q,
	}
}

func (r *wishlistRepository) GetWishlists(ctx context.Context, wishlistsIDs []string) ([]*entity.Wishlist, error) {
	query := `SELECT 
		id, user_id, name
		FROM wishlists 
		WHERE id = ANY($1)`
	rows, err := r.q.Query(ctx, query, wishlistsIDs)
	if err != nil {
		return nil, err
	}
	var wishlists []*entity.Wishlist
	err = pgxscan.ScanAll(&wishlists, rows)
	if err != nil {
		return nil, err
	}
	return wishlists, nil
}

func (r *wishlistRepository) GetWishlist(ctx context.Context, wishlistID string) (*entity.Wishlist, error) {
	query := `SELECT 
		id, user_id, name
		FROM wishlists 
		WHERE id = $1`
	var wishlist entity.Wishlist
	err := r.q.QueryRow(ctx, query, wishlistID).Scan(
		&wishlist.ID, &wishlist.UserID, &wishlist.Name,
	)
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

func (r *wishlistRepository) UpdateWishlist(ctx context.Context, wishlist *entity.Wishlist) error {
	query := `
		UPDATE wishlists 
		SET name = $2
		WHERE id = $1`
	ctag, err := r.q.Exec(ctx, query, wishlist.ID, wishlist.Name)
	if err != nil {
		return err
	}
	if ctag.RowsAffected() == 0 {
		return entity.ErrWishlistDoesNotExist
	}
	return nil
}

func (r *wishlistRepository) CreateWishlist(ctx context.Context, wishlist *entity.Wishlist) error {
	query := `INSERT INTO wishlists(id, user_id, name) VALUES ($1, $2, $3)`
	_, err := r.q.Exec(ctx, query, wishlist.ID, wishlist.UserID, wishlist.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishlistRepository) DeleteWishlist(ctx context.Context, wishlistID string) error {
	query := `DELETE FROM wishlists 
		WHERE id = $1`
	ctag, err := r.q.Exec(ctx, query, wishlistID)
	if err != nil {
		return err
	}
	if ctag.RowsAffected() == 0 {
		return entity.ErrWishlistDoesNotExist
	}
	return nil
}
