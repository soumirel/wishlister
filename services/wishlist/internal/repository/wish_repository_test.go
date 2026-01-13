package repository

import (
	"context"
	"testing"
	"time"

	entity "github.com/soumirel/wishlister/wishlist/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWishRepository_GetWish(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishRepository(q)
	ctx := context.Background()

	tests := []struct {
		name       string
		wishlistID string
		wishID     string
		wantErr    bool
		errType    error
		verify     func(t *testing.T, wish *entity.Wish)
	}{
		{
			name:       "should return existing wish",
			wishlistID: "alice_wishlist_1",
			wishID:     "alice_wish_1",
			wantErr:    false,
			verify: func(t *testing.T, wish *entity.Wish) {
				assert.Equal(t, "alice_wish_1", wish.ID)
				assert.Equal(t, "alice_wishlist_1", wish.WishlistID)
				assert.Equal(t, "Alice Wish 1", wish.Name)
			},
		},
		{
			name:       "should return error for non-existent wish",
			wishlistID: "alice_wishlist_1",
			wishID:     "non-existent-wish",
			wantErr:    true,
			errType:    entity.ErrWishDoesNotExist,
		},
		{
			name:       "should return error for non-existent wishlist",
			wishlistID: "non-existent-wishlist",
			wishID:     "alice_wish_1",
			wantErr:    true,
			errType:    entity.ErrWishDoesNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wish, err := repo.GetWish(ctx, tt.wishlistID, tt.wishID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.Nil(t, wish)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, wish)
			}
		})
	}
}

func TestWishRepository_GetWishesFromWishlist(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishRepository(q)
	ctx := context.Background()

	tests := []struct {
		name       string
		wishlistID string
		wantErr    bool
		verify     func(t *testing.T, wishes []*entity.Wish)
	}{
		{
			name:       "should return all wishes from wishlist",
			wishlistID: "alice_wishlist_1",
			wantErr:    false,
			verify: func(t *testing.T, wishes []*entity.Wish) {
				assert.GreaterOrEqual(t, len(wishes), 2, "should have at least 2 wishes")
				wishMap := make(map[string]*entity.Wish)
				for _, wish := range wishes {
					wishMap[wish.ID] = wish
					assert.Equal(t, "alice_wishlist_1", wish.WishlistID)
				}
				assert.Contains(t, wishMap, "alice_wish_1")
				assert.Contains(t, wishMap, "alice_wish_2")
			},
		},
		{
			name:       "should return empty list for wishlist with no wishes",
			wishlistID: "alice_wishlist_2",
			wantErr:    false,
			verify: func(t *testing.T, wishes []*entity.Wish) {
				assert.Empty(t, wishes)
			},
		},
		{
			name:       "should return empty list for non-existent wishlist",
			wishlistID: "non-existent-wishlist",
			wantErr:    false,
			verify: func(t *testing.T, wishes []*entity.Wish) {
				assert.Empty(t, wishes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wishes, err := repo.GetWishesFromWishlist(ctx, tt.wishlistID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, wishes)
			}
		})
	}
}

func TestWishRepository_CreateWish(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishRepository(q)
	ctx := context.Background()

	tests := []struct {
		name    string
		wish    *entity.Wish
		wantErr bool
		verify  func(t *testing.T, wishlistID, wishID string)
	}{
		{
			name: "should create wish without reservation",
			wish: &entity.Wish{
				ID:          "test-wish-1",
				WishlistID:  "alice_wishlist_1",
				Name:        "Test Wish 1",
				Reservation: nil,
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID, wishID string) {
				wish, err := repo.GetWish(ctx, wishlistID, wishID)
				require.NoError(t, err)
				assert.Equal(t, "test-wish-1", wish.ID)
				assert.Equal(t, "Test Wish 1", wish.Name)
				assert.Nil(t, wish.Reservation)
			},
		},
		{
			name: "should create wish with reservation",
			wish: &entity.Wish{
				ID:         "test-wish-2",
				WishlistID: "alice_wishlist_1",
				Name:       "Test Wish 2",
				Reservation: &entity.WishReservation{
					ID:               "test-reservation-1",
					ReservedByUserID: "alice",
					ReservedAt:       time.Now(),
				},
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID, wishID string) {
				wish, err := repo.GetWish(ctx, wishlistID, wishID)
				require.NoError(t, err)
				assert.Equal(t, "test-wish-2", wish.ID)
				assert.Equal(t, "Test Wish 2", wish.Name)
				require.NotNil(t, wish.Reservation)
				assert.Equal(t, "test-reservation-1", wish.Reservation.ID)
				assert.Equal(t, "alice", wish.Reservation.ReservedByUserID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateWish(ctx, tt.wish)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, tt.wish.WishlistID, tt.wish.ID)
			}
		})
	}
}

func TestWishRepository_UpdateWish(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishRepository(q)
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(t *testing.T) *entity.Wish
		wish    *entity.Wish
		wantErr bool
		verify  func(t *testing.T, wishlistID, wishID string)
	}{
		{
			name: "should update wish name",
			setup: func(t *testing.T) *entity.Wish {
				wish := &entity.Wish{
					ID:         "update-test-wish-1",
					WishlistID: "alice_wishlist_1",
					Name:       "Original Name",
				}
				err := repo.CreateWish(ctx, wish)
				require.NoError(t, err)
				return wish
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID, wishID string) {
				wish, err := repo.GetWish(ctx, wishlistID, wishID)
				require.NoError(t, err)
				assert.Equal(t, "Updated Name", wish.Name)
			},
		},
		{
			name: "should add reservation to wish",
			setup: func(t *testing.T) *entity.Wish {
				wish := &entity.Wish{
					ID:         "update-test-wish-2",
					WishlistID: "alice_wishlist_1",
					Name:       "Wish Without Reservation",
				}
				err := repo.CreateWish(ctx, wish)
				require.NoError(t, err)
				return wish
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID, wishID string) {
				wish, err := repo.GetWish(ctx, wishlistID, wishID)
				require.NoError(t, err)
				require.NotNil(t, wish.Reservation)
				assert.Equal(t, "test-reservation-2", wish.Reservation.ID)
				assert.Equal(t, "bob", wish.Reservation.ReservedByUserID)
			},
		},
		{
			name: "should remove reservation from wish",
			setup: func(t *testing.T) *entity.Wish {
				wish := &entity.Wish{
					ID:         "update-test-wish-3",
					WishlistID: "alice_wishlist_1",
					Name:       "Wish With Reservation",
					Reservation: &entity.WishReservation{
						ID:               "test-reservation-3",
						ReservedByUserID: "alice",
						ReservedAt:       time.Now(),
					},
				}
				err := repo.CreateWish(ctx, wish)
				require.NoError(t, err)
				return wish
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID, wishID string) {
				wish, err := repo.GetWish(ctx, wishlistID, wishID)
				require.NoError(t, err)
				assert.Nil(t, wish.Reservation)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalWish := tt.setup(t)
			wish := &entity.Wish{
				ID:         originalWish.ID,
				WishlistID: originalWish.WishlistID,
				Name:       originalWish.Name,
				Reservation: originalWish.Reservation,
			}

			// Modify wish based on test case
			switch tt.name {
			case "should update wish name":
				wish.Name = "Updated Name"
			case "should add reservation to wish":
				wish.Reservation = &entity.WishReservation{
					ID:               "test-reservation-2",
					ReservedByUserID: "bob",
					ReservedAt:       time.Now(),
				}
			case "should remove reservation from wish":
				wish.Reservation = nil
			}

			err := repo.UpdateWish(ctx, wish)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, wish.WishlistID, wish.ID)
			}
		})
	}
}

func TestWishRepository_DeleteWish(t *testing.T) {
	q := setupTestDB(t)
	repo := newWishRepository(q)
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(t *testing.T) (string, string)
		wishlistID string
		wishID     string
		wantErr    bool
		verify     func(t *testing.T, wishlistID, wishID string)
	}{
		{
			name: "should delete existing wish",
			setup: func(t *testing.T) (string, string) {
				wish := &entity.Wish{
					ID:         "delete-test-wish",
					WishlistID: "alice_wishlist_1",
					Name:       "Wish To Delete",
				}
				err := repo.CreateWish(ctx, wish)
				require.NoError(t, err)
				return wish.WishlistID, wish.ID
			},
			wantErr: false,
			verify: func(t *testing.T, wishlistID, wishID string) {
				_, err := repo.GetWish(ctx, wishlistID, wishID)
				assert.Error(t, err)
				assert.ErrorIs(t, err, entity.ErrWishDoesNotExist)
			},
		},
		{
			name:       "should delete non-existent wish without error",
			wishlistID: "alice_wishlist_1",
			wishID:     "non-existent-wish",
			wantErr:    false,
			verify: func(t *testing.T, wishlistID, wishID string) {
				_, err := repo.GetWish(ctx, wishlistID, wishID)
				assert.Error(t, err)
				assert.ErrorIs(t, err, entity.ErrWishDoesNotExist)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wishlistID := tt.wishlistID
			wishID := tt.wishID
			if tt.setup != nil {
				wishlistID, wishID = tt.setup(t)
			}

			err := repo.DeleteWish(ctx, wishlistID, wishID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, wishlistID, wishID)
			}
		})
	}
}
