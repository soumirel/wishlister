package repository

import (
	"context"
	"testing"

	entity "github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWishlistRepository_GetWishlist(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishlistRepository(q)
	ctx := context.Background()

	tests := []struct {
		name       string
		wishlistID string
		wantErr    bool
		verify     func(t *testing.T, wishlist *entity.Wishlist)
	}{
		{
			name:       "should return existing wishlist",
			wishlistID: "alice_wishlist_1",
			wantErr:    false,
			verify: func(t *testing.T, wishlist *entity.Wishlist) {
				assert.Equal(t, "alice_wishlist_1", wishlist.ID)
				assert.Equal(t, "alice", wishlist.UserID)
				assert.Equal(t, "Alice Wishlist 1", wishlist.Name)
			},
		},
		{
			name:       "should return existing wishlist bob",
			wishlistID: "bob_wishlist_1",
			wantErr:    false,
			verify: func(t *testing.T, wishlist *entity.Wishlist) {
				assert.Equal(t, "bob_wishlist_1", wishlist.ID)
				assert.Equal(t, "bob", wishlist.UserID)
				assert.Equal(t, "Bob Wishlist 1", wishlist.Name)
			},
		},
		{
			name:       "should return error for non-existent wishlist",
			wishlistID: "non-existent-wishlist",
			wantErr:    true,
			verify: func(t *testing.T, wishlist *entity.Wishlist) {
				// This should not be called for error cases
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wishlist, err := repo.GetWishlist(ctx, tt.wishlistID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, wishlist)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, wishlist)
			}
		})
	}
}

func TestWishlistRepository_GetWishlists(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishlistRepository(q)
	ctx := context.Background()

	tests := []struct {
		name         string
		wishlistsIDs []string
		wantErr      bool
		verify       func(t *testing.T, wishlists []*entity.Wishlist)
	}{
		{
			name:         "should return multiple wishlists",
			wishlistsIDs: []string{"alice_wishlist_1", "alice_wishlist_2", "bob_wishlist_1"},
			wantErr:      false,
			verify: func(t *testing.T, wishlists []*entity.Wishlist) {
				assert.Equal(t, 3, len(wishlists))
				wishlistMap := make(map[string]*entity.Wishlist)
				for _, wl := range wishlists {
					wishlistMap[wl.ID] = wl
				}
				assert.Contains(t, wishlistMap, "alice_wishlist_1")
				assert.Contains(t, wishlistMap, "alice_wishlist_2")
				assert.Contains(t, wishlistMap, "bob_wishlist_1")
			},
		},
		{
			name:         "should return single wishlist",
			wishlistsIDs: []string{"alice_wishlist_1"},
			wantErr:      false,
			verify: func(t *testing.T, wishlists []*entity.Wishlist) {
				assert.Equal(t, 1, len(wishlists))
				assert.Equal(t, "alice_wishlist_1", wishlists[0].ID)
			},
		},
		{
			name:         "should return empty list for non-existent IDs",
			wishlistsIDs: []string{"non-existent-1", "non-existent-2"},
			wantErr:      false,
			verify: func(t *testing.T, wishlists []*entity.Wishlist) {
				assert.Empty(t, wishlists)
			},
		},
		{
			name:         "should return empty list for empty input",
			wishlistsIDs: []string{},
			wantErr:      false,
			verify: func(t *testing.T, wishlists []*entity.Wishlist) {
				assert.Empty(t, wishlists)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wishlists, err := repo.GetWishlists(ctx, tt.wishlistsIDs)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, wishlists)
			}
		})
	}
}

func TestWishlistRepository_CreateWishlist(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishlistRepository(q)
	ctx := context.Background()

	tests := []struct {
		name     string
		wishlist *entity.Wishlist
		wantErr  bool
		verify   func(t *testing.T, wishlistID string)
	}{
		{
			name: "should create new wishlist",
			wishlist: &entity.Wishlist{
				ID:     "test-wishlist-1",
				UserID: "alice",
				Name:   "Test Wishlist 1",
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID string) {
				wishlist, err := repo.GetWishlist(ctx, wishlistID)
				require.NoError(t, err)
				assert.Equal(t, "test-wishlist-1", wishlist.ID)
				assert.Equal(t, "alice", wishlist.UserID)
				assert.Equal(t, "Test Wishlist 1", wishlist.Name)
			},
		},
		{
			name: "should create another new wishlist",
			wishlist: &entity.Wishlist{
				ID:     "test-wishlist-2",
				UserID: "bob",
				Name:   "Test Wishlist 2",
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID string) {
				wishlist, err := repo.GetWishlist(ctx, wishlistID)
				require.NoError(t, err)
				assert.Equal(t, "test-wishlist-2", wishlist.ID)
				assert.Equal(t, "bob", wishlist.UserID)
				assert.Equal(t, "Test Wishlist 2", wishlist.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateWishlist(ctx, tt.wishlist)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, tt.wishlist.ID)
			}
		})
	}
}

func TestWishlistRepository_UpdateWishlist(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishlistRepository(q)
	ctx := context.Background()

	tests := []struct {
		name     string
		setup    func(t *testing.T) *entity.Wishlist
		wishlist *entity.Wishlist
		wantErr  bool
		errType  error
		verify   func(t *testing.T, wishlistID string)
	}{
		{
			name: "should update existing wishlist",
			setup: func(t *testing.T) *entity.Wishlist {
				wishlist := &entity.Wishlist{
					ID:     "update-test-wishlist",
					UserID: "alice",
					Name:   "Original Name",
				}
				err := repo.CreateWishlist(ctx, wishlist)
				require.NoError(t, err)
				return wishlist
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID string) {
				wishlist, err := repo.GetWishlist(ctx, wishlistID)
				require.NoError(t, err)
				assert.Equal(t, "Updated Name", wishlist.Name)
			},
		},
		{
			name: "should return error for non-existent wishlist",
			wishlist: &entity.Wishlist{
				ID:     "non-existent-wishlist",
				UserID: "alice",
				Name:   "Updated Name",
			},
			wantErr: true,
			errType: entity.ErrWishlistDoesNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wishlist := tt.wishlist
			if tt.setup != nil {
				wishlist = tt.setup(t)
				wishlist.Name = "Updated Name"
			}

			err := repo.UpdateWishlist(ctx, wishlist)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, wishlist.ID)
			}
		})
	}
}

func TestWishlistRepository_DeleteWishlist(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishlistRepository(q)
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wishlistID string
		wantErr    bool
		errType    error
		verify     func(t *testing.T, wishlistID string)
	}{
		{
			name: "should delete existing wishlist",
			setup: func(t *testing.T) string {
				wishlist := &entity.Wishlist{
					ID:     "delete-test-wishlist",
					UserID: "alice",
					Name:   "Wishlist To Delete",
				}
				err := repo.CreateWishlist(ctx, wishlist)
				require.NoError(t, err)
				return wishlist.ID
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID string) {
				_, err := repo.GetWishlist(ctx, wishlistID)
				assert.Error(t, err)
			},
		},
		{
			name:       "should return error for non-existent wishlist",
			wishlistID: "non-existent-wishlist",
			wantErr:    true,
			errType:    entity.ErrWishlistDoesNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wishlistID := tt.wishlistID
			if tt.setup != nil {
				wishlistID = tt.setup(t)
			}

			err := repo.DeleteWishlist(ctx, wishlistID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, wishlistID)
			}
		})
	}
}
