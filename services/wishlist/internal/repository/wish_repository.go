package repository

import (
	"context"
	"errors"

	"github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
	"github.com/soumirel/wishlister/services/wishlist/internal/domain/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type wishRepository struct {
	q Querier
}

func newWishRepository(q Querier) repository.WishRepository {
	return &wishRepository{
		q: q,
	}
}

// scanWishRow scans pgx.Row to entity.Wish
// rows can be obtained from different sources (Query, QueryRow),
// so the calling code must check ErrNoRows itself
func (r *wishRepository) scanWishRow(row pgx.Row) (*entity.Wish, error) {
	var (
		wish        entity.Wish
		reservation struct {
			ID               pgtype.Text
			ReservedByUserID pgtype.Text
			ReservedAt       pgtype.Timestamptz
		}
	)
	err := row.Scan(
		&wish.ID, &wish.WishlistID, &wish.Name,
		&reservation.ID, &reservation.ReservedByUserID, &reservation.ReservedAt,
	)
	if err != nil {
		return nil, err
	}
	if reservation.ID.Valid {
		wish.Reservation = &entity.WishReservation{
			ID:               reservation.ID.String,
			ReservedByUserID: reservation.ReservedByUserID.String,
			ReservedAt:       reservation.ReservedAt.Time,
		}
	}
	return &wish, nil
}

func (r *wishRepository) GetWish(ctx context.Context, wishlistID, wishID string) (*entity.Wish, error) {
	query := `
		SELECT 
			w.id, w.wishlist_id, w.name,
			r.id, r.reserved_by_user_id, r.reserved_at
		FROM wishes w
		LEFT JOIN wish_reservations r
			ON r.wish_id = w.id
		WHERE w.wishlist_id = $1
			AND w.id = $2`
	wish, err := r.scanWishRow(r.q.QueryRow(ctx, query, wishlistID, wishID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrWishDoesNotExist
		}
		return nil, err
	}
	return wish, nil
}

func (r *wishRepository) GetWishesFromWishlist(ctx context.Context, wishlistID string) ([]*entity.Wish, error) {
	query := `
		SELECT 
			w.id, w.wishlist_id, w.name,
			r.id, r.reserved_by_user_id, r.reserved_at
		FROM wishes w
		LEFT JOIN wish_reservations r
			ON r.wish_id = w.id
		WHERE w.wishlist_id = $1`

	rows, err := r.q.Query(ctx, query, wishlistID)
	wishes, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*entity.Wish, error) {
		wish, err := r.scanWishRow(row)
		if err != nil {
			return nil, err
		}
		return wish, nil
	})
	if err != nil {
		return nil, err
	}
	return wishes, nil
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

	var (
		reservationQuery string
		reservationArgs  []any
	)
	if wish.Reservation == nil {
		reservationQuery = `
			DELETE FROM wish_reservations
			WHERE wish_id = $1`
		reservationArgs = []any{wish.ID}
	} else {
		reservation := wish.Reservation
		reservationQuery = `
			INSERT INTO wish_reservations (
				id, wish_id, reserved_by_user_id, reserved_at
			) VALUES ($1, $2, $3, $4)
				ON CONFLICT (id) 
				DO NOTHING`
		reservationArgs = []any{
			reservation.ID,
			wish.ID,
			reservation.ReservedByUserID,
			reservation.ReservedAt,
		}
	}

	_, err = r.q.Exec(ctx, reservationQuery, reservationArgs...)
	if err != nil {
		return err
	}
	return nil
}

func (r *wishRepository) CreateWish(ctx context.Context, wish *entity.Wish) error {
	insertWishQuery := `
			INSERT INTO wishes(id, wishlist_id, name)
			VALUES ($1, $2, $3)`
	_, err := r.q.Exec(ctx, insertWishQuery,
		wish.ID, wish.WishlistID, wish.Name,
	)
	if err != nil {
		return err
	}

	if wish.Reservation == nil {
		return nil
	}
	reservation := wish.Reservation
	insertReservationQuery := `
			INSERT INTO wish_reservations(
				id, wish_id, reserved_by_user_id, reserved_at
			) VALUES($1, $2, $3, $4)`
	_, err = r.q.Exec(ctx, insertReservationQuery,
		reservation.ID, wish.ID, reservation.ReservedByUserID, reservation.ReservedAt,
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
