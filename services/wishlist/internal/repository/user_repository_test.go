package repository

import (
	"context"
	"testing"

	entity "github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_GetUsers(t *testing.T) {
	q := setupTestDB(t)
	repo := newUserRepository(q)
	ctx := context.Background()

	tests := []struct {
		name    string
		wantErr bool
		verify  func(t *testing.T, users []*entity.User)
	}{
		{
			name:    "should return all users",
			wantErr: false,
			verify: func(t *testing.T, users []*entity.User) {
				assert.GreaterOrEqual(t, len(users), 2, "should have at least 2 users")
				userMap := make(map[string]*entity.User)
				for _, user := range users {
					userMap[user.ID] = user
				}
				assert.Contains(t, userMap, "alice", "should contain alice")
				assert.Contains(t, userMap, "bob", "should contain bob")
				if alice, ok := userMap["alice"]; ok {
					assert.Equal(t, "Alice", alice.Name)
				}
				if bob, ok := userMap["bob"]; ok {
					assert.Equal(t, "Bob", bob.Name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := repo.GetUsers(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.verify(t, users)
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	q := setupTestDB(t)
	repo := newUserRepository(q)
	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		wantErr bool
		want    *entity.User
		errType error
	}{
		{
			name:    "should return existing user",
			userID:  "alice",
			wantErr: false,
			want: &entity.User{
				ID:   "alice",
				Name: "Alice",
			},
		},
		{
			name:    "should return existing user bob",
			userID:  "bob",
			wantErr: false,
			want: &entity.User{
				ID:   "bob",
				Name: "Bob",
			},
		},
		{
			name:    "should return error for non-existent user",
			userID:  "non-existent",
			wantErr: true,
			errType: entity.ErrUserDoesNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetUser(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.Nil(t, user)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.ID, user.ID)
			assert.Equal(t, tt.want.Name, user.Name)
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	q := setupTestDB(t)
	repo := newUserRepository(q)
	ctx := context.Background()

	tests := []struct {
		name    string
		user    *entity.User
		wantErr bool
		verify  func(t *testing.T, userID string)
	}{
		{
			name: "should create new user",
			user: &entity.User{
				ID:   "test-user-1",
				Name: "Test User 1",
			},
			wantErr: false,
			verify: func(t *testing.T, userID string) {
				user, err := repo.GetUser(ctx, userID)
				require.NoError(t, err)
				assert.Equal(t, "test-user-1", user.ID)
				assert.Equal(t, "Test User 1", user.Name)
			},
		},
		{
			name: "should create another new user",
			user: &entity.User{
				ID:   "test-user-2",
				Name: "Test User 2",
			},
			wantErr: false,
			verify: func(t *testing.T, userID string) {
				user, err := repo.GetUser(ctx, userID)
				require.NoError(t, err)
				assert.Equal(t, "test-user-2", user.ID)
				assert.Equal(t, "Test User 2", user.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateUser(ctx, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, tt.user.ID)
			}
		})
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	q := setupTestDB(t)
	repo := newUserRepository(q)
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		userID  string
		wantErr bool
		verify  func(t *testing.T, userID string)
	}{
		{
			name: "should delete existing user",
			setup: func(t *testing.T) string {
				user := &entity.User{
					ID:   "delete-test-user",
					Name: "Delete Test User",
				}
				err := repo.CreateUser(ctx, user)
				require.NoError(t, err)
				return user.ID
			},
			wantErr: false,
			verify: func(t *testing.T, userID string) {
				_, err := repo.GetUser(ctx, userID)
				assert.Error(t, err)
				assert.ErrorIs(t, err, entity.ErrUserDoesNotExist)
			},
		},
		{
			name:    "should delete non-existent user without error",
			userID:  "non-existent-delete",
			wantErr: false,
			verify: func(t *testing.T, userID string) {
				_, err := repo.GetUser(ctx, userID)
				assert.Error(t, err)
				assert.ErrorIs(t, err, entity.ErrUserDoesNotExist)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.userID
			if tt.setup != nil {
				userID = tt.setup(t)
			}

			err := repo.DeleteUser(ctx, userID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, userID)
			}
		})
	}
}
