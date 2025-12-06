package repository

// import (
// 	"context"
// 	"wishlister/internal/domain"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type reservationRepository struct {
// 	db *pgxpool.Pool
// }

// func NewReservationRepository(pool *pgxpool.Pool) *reservationRepository {
// 	return &reservationRepository{
// 		db: pool,
// 	}
// }

// func (r *reservationRepository) Reserve(ctx context.Context, reservation domain.Reservation) error {
// 	query := `INSERT INTO reservations(wish_id, user_id) VALUES ($1, $2)`
// 	_, err := r.db.Exec(ctx, query, reservation.WishID, reservation.UserID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *reservationRepository) Release(ctx context.Context, reservation domain.Reservation) error {
// 	query := `DELETE FROM reservations WHERE wish_id = $1,user_id = $2`
// 	_, err := r.db.Exec(ctx, query, reservation.WishID, reservation.UserID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
