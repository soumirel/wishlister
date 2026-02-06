package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	entity "github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
)

type userIdentityRepository struct {
	q Querier
}

func newUserIdentityRepository(q Querier) *userIdentityRepository {
	return &userIdentityRepository{
		q: q,
	}
}

func (r *userIdentityRepository) GetUserIdByExternalIdentity(ctx context.Context, externalIdentity entity.ExternalIdentity) (string, error) {
	query := `
		SELECT user_id 
		FROM external_user_identities
		WHERE external_id = $1
			AND provider = $2`
	var (
		userID string
	)
	err := r.q.QueryRow(ctx, query, externalIdentity.ExternalID, externalIdentity.Provider).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", entity.ErrUserIdentityDoesNotExist
		}
		return "", err
	}
	return userID, nil
}

func (r *userIdentityRepository) SaveIdentity(ctx context.Context, userIdentity *entity.UserIdentity) error {
	query := `
		INSERT INTO external_user_identities(id, user_id, external_id, provider)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING`
	ctag, err := r.q.Exec(ctx, query,
		userIdentity.ID,
		userIdentity.UserID,
		userIdentity.ExternalIdentity.ExternalID,
		userIdentity.ExternalIdentity.Provider,
	)
	if err != nil {
		return err
	}
	if ctag.RowsAffected() == 0 {
		return entity.ErrExternalUserIdentityConflict
	}
	return nil
}
