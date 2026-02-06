package repository

import (
	"context"
	"testing"

	entity "github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserIdentityRepository_GetUserIdByExternalIdentity(t *testing.T) {
	q := setupTestDB(t)
	userRepo := newUserRepository(q)
	repo := newUserIdentityRepository(q)
	ctx := context.Background()

	// Setup: create a user and user identity
	testUserID := "test-identity-user"
	testExternalID := "telegram-12345"
	testProvider := entity.IdentityProvider("telegram")

	// Create user first
	user := &entity.User{
		ID:   testUserID,
		Name: "Test Identity User",
	}
	err := userRepo.CreateUser(ctx, user)
	require.NoError(t, err)

	userIdentity := entity.NewUserIdentity(testUserID, entity.ExternalIdentity{
		ExternalID: testExternalID,
		Provider:   testProvider,
	})
	err = repo.SaveIdentity(ctx, userIdentity)
	require.NoError(t, err)

	tests := []struct {
		name       string
		externalID string
		provider   entity.IdentityProvider
		wantErr    bool
		wantUserID string
		errType    error
	}{
		{
			name:       "should return user ID for existing identity",
			externalID: testExternalID,
			provider:   testProvider,
			wantErr:    false,
			wantUserID: testUserID,
		},
		{
			name:       "should return error for non-existent external ID",
			externalID: "non-existent-id",
			provider:   testProvider,
			wantErr:    true,
			errType:    entity.ErrUserIdentityDoesNotExist,
		},
		{
			name:       "should return error for non-existent provider",
			externalID: testExternalID,
			provider:   entity.IdentityProvider("unknown"),
			wantErr:    true,
			errType:    entity.ErrUserIdentityDoesNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			externalIdentity := entity.ExternalIdentity{
				ExternalID: tt.externalID,
				Provider:   tt.provider,
			}
			userID, err := repo.GetUserIdByExternalIdentity(ctx, externalIdentity)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.Empty(t, userID)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantUserID, userID)
		})
	}
}

func TestUserIdentityRepository_SaveIdentity(t *testing.T) {
	q := setupTestDB(t)
	userRepo := newUserRepository(q)
	repo := newUserIdentityRepository(q)
	ctx := context.Background()

	tests := []struct {
		name       string
		userID     string
		externalID string
		provider   entity.IdentityProvider
		setup      func(t *testing.T, userID string)
		wantErr    bool
		errType    error
		verify     func(t *testing.T, userID string, externalID string, provider entity.IdentityProvider)
	}{
		{
			name:       "should save new identity",
			userID:     "save-test-user-1",
			externalID: "telegram-11111",
			provider:   entity.IdentityProvider("telegram"),
			setup: func(t *testing.T, userID string) {
				user := &entity.User{
					ID:   userID,
					Name: "Save Test User 1",
				}
				err := userRepo.CreateUser(ctx, user)
				require.NoError(t, err)
			},
			wantErr: false,
			verify: func(t *testing.T, userID string, externalID string, provider entity.IdentityProvider) {
				externalIdentity := entity.ExternalIdentity{
					ExternalID: externalID,
					Provider:   provider,
				}
				foundUserID, err := repo.GetUserIdByExternalIdentity(ctx, externalIdentity)
				require.NoError(t, err)
				assert.Equal(t, userID, foundUserID)
			},
		},
		{
			name:       "should return error on conflict",
			userID:     "save-test-user-2",
			externalID: "telegram-22222",
			provider:   entity.IdentityProvider("telegram"),
			setup: func(t *testing.T, userID string) {
				user := &entity.User{
					ID:   userID,
					Name: "Save Test User 2",
				}
				err := userRepo.CreateUser(ctx, user)
				require.NoError(t, err)
			},
			wantErr: false,
			verify: func(t *testing.T, userID string, externalID string, provider entity.IdentityProvider) {
				// Try to save the same identity again with different user
				duplicateUserID := "save-test-user-3"
				duplicateUser := &entity.User{
					ID:   duplicateUserID,
					Name: "Save Test User 3",
				}
				err := userRepo.CreateUser(ctx, duplicateUser)
				require.NoError(t, err)

				duplicateIdentity := entity.NewUserIdentity(duplicateUserID, entity.ExternalIdentity{
					ExternalID: externalID,
					Provider:   provider,
				})
				err = repo.SaveIdentity(ctx, duplicateIdentity)
				assert.Error(t, err)
				assert.ErrorIs(t, err, entity.ErrExternalUserIdentityConflict)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, tt.userID)
			}

			userIdentity := entity.NewUserIdentity(tt.userID, entity.ExternalIdentity{
				ExternalID: tt.externalID,
				Provider:   tt.provider,
			})
			err := repo.SaveIdentity(ctx, userIdentity)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, tt.userID, tt.externalID, tt.provider)
			}
		})
	}
}
