package repository

import (
	"context"
	"wishlister/internal/domain"
	"wishlister/internal/pkg"
	"wishlister/internal/service/wishlist"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type wishlistRepository struct {
	db pkg.DbExecutor
}

func NewWishlistRepository(db pkg.DbExecutor) *wishlistRepository {
	return &wishlistRepository{
		db: db,
	}
}

func (r *wishlistRepository) GetUserWishlists(ctx context.Context, userID string) ([]*domain.Wishlist, error) {
	query := `SELECT 
		id, user_id, name
		FROM wishlists 
		WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	var wishlists []*domain.Wishlist
	err = pgxscan.ScanAll(&wishlists, rows)
	if err != nil {
		return nil, err
	}
	return wishlists, nil
}

func (r *wishlistRepository) GetWishlists(ctx context.Context, filter wishlist.WishlistFilter) ([]*domain.Wishlist, error) {
	query := `SELECT 
		id, user_id, name
		FROM wishlists 
		WHERE id = ANY($1)`
	rows, err := r.db.Query(ctx, query, filter.WishlistIDs)
	if err != nil {
		return nil, err
	}
	var wishlists []*domain.Wishlist
	err = pgxscan.ScanAll(&wishlists, rows)
	if err != nil {
		return nil, err
	}
	return wishlists, nil
}

func (r *wishlistRepository) GetWishlist(ctx context.Context, wishlistID string) (*domain.Wishlist, error) {
	query := `SELECT 
			wl.id as wishlist_id, wl.user_id as wishlist_user, wl.name as wishlist_name,
			w.id as wish_id, w.name as wish_name
		FROM wishlists wl
		LEFT JOIN wishes w
			ON w.wishlist_id = wl.id
		WHERE wl.id = $1`
	rows, err := r.db.Query(ctx, query, wishlistID)
	if err != nil {
		return &domain.Wishlist{}, err
	}
	var (
		wishlist               *domain.Wishlist
		wlID, wlUserID, wlName string
		wID, wName             pgtype.Text
	)
	pgx.ForEachRow(rows,
		[]any{&wlID, &wlUserID, &wlName, &wID, &wName},
		func() error {
			if wishlist == nil {
				wishlist = &domain.Wishlist{
					ID:     wlID,
					UserID: wlUserID,
					Name:   wlName,
					Wishes: make(map[string]*domain.Wish),
				}
			}

			if wID.Valid {
				wishlist.Wishes[wID.String] = &domain.Wish{
					ID:         wID.String,
					UserID:     wlUserID,
					WishlistID: wlID,
					Name:       wName.String,
				}
			}
			return nil
		})
	if wishlist == nil {
		return nil, domain.ErrWishlistDoesNotExist
	}
	return wishlist, nil
}

func (r *wishlistRepository) UpdateWishlist(ctx context.Context, wishlist *domain.Wishlist) error {
	query := `
		UPDATE wishlists 
		SET name = $3
		WHERE user_id = $1
			AND id = $2`
	ctag, err := r.db.Exec(ctx, query, wishlist.UserID, wishlist.ID, wishlist.Name)
	if err != nil {
		return err
	}
	if ctag.RowsAffected() == 0 {
		return domain.ErrWishlistDoesNotExist
	}
	return nil
}

func (r *wishlistRepository) CreateWishlist(ctx context.Context, wishlist *domain.Wishlist) error {
	query := `INSERT INTO wishlists(id, user_id, name) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, wishlist.ID, wishlist.UserID, wishlist.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishlistRepository) DeleteWishlist(ctx context.Context, wishlistID string) error {
	query := `DELETE FROM wishlists 
		WHERE id = $1`
	ctag, err := r.db.Exec(ctx, query, wishlistID)
	if err != nil {
		return err
	}
	if ctag.RowsAffected() == 0 {
		return domain.ErrWishlistDoesNotExist
	}
	return nil
}
