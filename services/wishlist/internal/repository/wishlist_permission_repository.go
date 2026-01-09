package repository

import (
	"context"
	"errors"

	"github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/wishlist/internal/domain/repository"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type wishlistPermissionRepo struct {
	q Querier
}

func newWishlistPersmissionRepository(q Querier) repository.WishlistPermissionRepository {
	return &wishlistPermissionRepo{
		q: q,
	}
}

func (r *wishlistPermissionRepo) GetPermissionToWishlist(ctx context.Context, userID, wishlistID string) (*entity.WishlistPermission, error) {
	query := `
			SELECT id, user_id, wishlist_id, permission_level
			FROM wishlist_permissions
			WHERE user_id = $1
				AND wishlist_id = $2`
	var permission entity.WishlistPermission
	err := r.q.QueryRow(ctx, query, userID, wishlistID).Scan(
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

func (r *wishlistPermissionRepo) GetPermissionsToWishlists(ctx context.Context, userID string) (entity.WishlistsPermissions, error) {
	query := `
			SELECT id, user_id, wishlist_id, permission_level
			FROM wishlist_permissions
			WHERE user_id = $1`
	var permissions []*entity.WishlistPermission
	rows, err := r.q.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	err = pgxscan.ScanAll(&permissions, rows)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *wishlistPermissionRepo) SaveWishlistPermission(ctx context.Context, permission *entity.WishlistPermission) error {
	query := `
			INSERT INTO wishlist_permissions(user_id, wishlist_id, permission_level)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, wishlist_id)
			DO UPDATE
				SET permission_level = excluded.permission_level`
	_, err := r.q.Exec(ctx, query, permission.UserID, permission.WishlistID, permission.Level)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishlistPermissionRepo) DeleteWishlistPermission(ctx context.Context, userID, wishlistID string) error {
	query := `
			DELETE FROM wishlist_permissions
			WHERE user_id = $1
				AND wishlist_id = $2`
	_, err := r.q.Exec(ctx, query, userID, wishlistID)
	if err != nil {
		return nil
	}
	return nil
}
