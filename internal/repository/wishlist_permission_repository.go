package repository

import (
	"context"
	"errors"
	"wishlister/internal/domain"
	"wishlister/internal/pkg"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type wishlistPermissionRepo struct {
	db pkg.DbExecutor
}

func NewWishlistPersmissionRepository(db pkg.DbExecutor) *wishlistPermissionRepo {
	return &wishlistPermissionRepo{
		db: db,
	}
}

func (r *wishlistPermissionRepo) GetPermissionToWishlist(ctx context.Context, userID, wishlistID string) (*domain.WishlistPermission, error) {
	query := `
			SELECT id, user_id, wishlist_id, permission_level
			FROM wishlists_permissions
			WHERE user_id = $1
				AND wishlist_id = $2`
	var permission domain.WishlistPermission
	err := r.db.QueryRow(ctx, query, userID, wishlistID).Scan(
		&permission.ID, &permission.UserID, &permission.WishlistID, &permission.Level,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

func (r *wishlistPermissionRepo) GetUserPermissionsToWishlists(ctx context.Context, userID string) (domain.WishlistsPersmmissions, error) {
	query := `
			SELECT id, user_id, wishlist_id, permission_level
			FROM wishlists_permissions
			WHERE user_id = $1`
	var permissions []*domain.WishlistPermission
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	err = pgxscan.ScanAll(&permissions, rows)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *wishlistPermissionRepo) SaveWishlistPermission(ctx context.Context, permission *domain.WishlistPermission) error {
	query := `
			INSERT INTO wishlists_permissions(user_id, wishlist_id, permission_level)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, wishlist_id)
			DO UPDATE
				SET permission_level = excluded.permission_level`
	_, err := r.db.Exec(ctx, query, permission.UserID, permission.WishlistID, permission.Level)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishlistPermissionRepo) DeleteWishlistPermission(ctx context.Context, userID, wishlistID string) error {
	query := `
			DELETE FROM wishlists_permissions
			WHERE user_id = $1
				AND wishlist_id = $2`
	_, err := r.db.Exec(ctx, query, userID, wishlistID)
	if err != nil {
		return nil
	}
	return nil
}
