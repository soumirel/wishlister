package repository

import (
	"context"
	"testing"

	entity "github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWishlistPermissionRepository_GetPermissionToWishlist(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishlistPersmissionRepository(q)
	ctx := context.Background()

	tests := []struct {
		name       string
		userID     string
		wishlistID string
		wantErr    bool
		verify     func(t *testing.T, permission *entity.WishlistPermission)
	}{
		{
			name:       "should return existing permission",
			userID:     "alice",
			wishlistID: "alice_wishlist_1",
			wantErr:    false,
			verify: func(t *testing.T, permission *entity.WishlistPermission) {
				require.NotNil(t, permission)
				assert.Equal(t, "alice", permission.UserID)
				assert.Equal(t, "alice_wishlist_1", permission.WishlistID)
				assert.Equal(t, entity.OwnerWishlistPermissionLevel, permission.Level)
			},
		},
		{
			name:       "should return nil for non-existent permission",
			userID:     "alice",
			wishlistID: "bob_wishlist_1",
			wantErr:    false,
			verify: func(t *testing.T, permission *entity.WishlistPermission) {
				assert.Nil(t, permission)
			},
		},
		{
			name:       "should return nil for non-existent user",
			userID:     "non-existent-user",
			wishlistID: "alice_wishlist_1",
			wantErr:    false,
			verify: func(t *testing.T, permission *entity.WishlistPermission) {
				assert.Nil(t, permission)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permission, err := repo.GetPermissionToWishlist(ctx, tt.userID, tt.wishlistID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, permission)
			}
		})
	}
}

func TestWishlistPermissionRepository_GetPermissionsToWishlists(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishlistPersmissionRepository(q)
	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		wantErr bool
		verify  func(t *testing.T, permissions entity.WishlistsPermissions)
	}{
		{
			name:    "should return all permissions for user",
			userID:  "alice",
			wantErr: false,
			verify: func(t *testing.T, permissions entity.WishlistsPermissions) {
				assert.GreaterOrEqual(t, len(permissions), 2, "should have at least 2 permissions")
				permissionMap := make(map[string]*entity.WishlistPermission)
				for _, perm := range permissions {
					permissionMap[perm.WishlistID] = perm
					assert.Equal(t, "alice", perm.UserID)
				}
				assert.Contains(t, permissionMap, "alice_wishlist_1")
				assert.Contains(t, permissionMap, "alice_wishlist_2")
			},
		},
		{
			name:    "should return empty list for user with no permissions",
			userID:  "non-existent-user",
			wantErr: false,
			verify: func(t *testing.T, permissions entity.WishlistsPermissions) {
				assert.Empty(t, permissions)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permissions, err := repo.GetPermissionsToWishlists(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, permissions)
			}
		})
	}
}

func TestWishlistPermissionRepository_SaveWishlistPermission(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishlistPersmissionRepository(q)
	ctx := context.Background()

	tests := []struct {
		name       string
		permission *entity.WishlistPermission
		wantErr    bool
		verify     func(t *testing.T, userID, wishlistID string)
	}{
		{
			name: "should save new permission",
			permission: &entity.WishlistPermission{
				UserID:     "alice",
				WishlistID: "bob_wishlist_1",
				Level:      entity.GuestWishlistPersmissionLevel,
			},
			wantErr: false,
			verify: func(t *testing.T, userID, wishlistID string) {
				permission, err := repo.GetPermissionToWishlist(ctx, userID, wishlistID)
				require.NoError(t, err)
				require.NotNil(t, permission)
				assert.Equal(t, userID, permission.UserID)
				assert.Equal(t, wishlistID, permission.WishlistID)
				assert.Equal(t, entity.GuestWishlistPersmissionLevel, permission.Level)
			},
		},
		{
			name: "should update existing permission",
			permission: &entity.WishlistPermission{
				UserID:     "alice",
				WishlistID: "alice_wishlist_1",
				Level:      entity.GuestWishlistPersmissionLevel,
			},
			wantErr: false,
			verify: func(t *testing.T, userID, wishlistID string) {
				permission, err := repo.GetPermissionToWishlist(ctx, userID, wishlistID)
				require.NoError(t, err)
				require.NotNil(t, permission)
				assert.Equal(t, entity.GuestWishlistPersmissionLevel, permission.Level, "permission level should be updated")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SaveWishlistPermission(ctx, tt.permission)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, tt.permission.UserID, tt.permission.WishlistID)
			}
		})
	}
}

func TestWishlistPermissionRepository_DeleteWishlistPermission(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishlistPersmissionRepository(q)
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(t *testing.T) (string, string)
		userID     string
		wishlistID string
		wantErr    bool
		verify     func(t *testing.T, userID, wishlistID string)
	}{
		{
			name: "should delete existing permission",
			setup: func(t *testing.T) (string, string) {
				permission := &entity.WishlistPermission{
					UserID:     "alice",
					WishlistID: "bob_wishlist_2",
					Level:      entity.GuestWishlistPersmissionLevel,
				}
				err := repo.SaveWishlistPermission(ctx, permission)
				require.NoError(t, err)
				return permission.UserID, permission.WishlistID
			},
			wantErr: false,
			verify: func(t *testing.T, userID, wishlistID string) {
				permission, err := repo.GetPermissionToWishlist(ctx, userID, wishlistID)
				require.NoError(t, err)
				assert.Nil(t, permission, "permission should be deleted")
			},
		},
		{
			name:       "should delete non-existent permission without error",
			userID:     "non-existent-user",
			wishlistID: "non-existent-wishlist",
			wantErr:    false,
			verify: func(t *testing.T, userID, wishlistID string) {
				permission, err := repo.GetPermissionToWishlist(ctx, userID, wishlistID)
				require.NoError(t, err)
				assert.Nil(t, permission)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.userID
			wishlistID := tt.wishlistID
			if tt.setup != nil {
				userID, wishlistID = tt.setup(t)
			}

			err := repo.DeleteWishlistPermission(ctx, userID, wishlistID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, userID, wishlistID)
			}
		})
	}
}
