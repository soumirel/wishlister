package entity

import (
	"errors"
)

type WishlistAction string

const (
	ReadWishlistAction        WishlistAction = "read"
	ModifyWishlistAction      WishlistAction = "modify"
	ReserveWishWishlistAction                = "reserve"
)

type WishlistPermissionLevel string

const (
	OwnerWishlistPermissionLevel  WishlistPermissionLevel = "owner"
	GuestWishlistPersmissionLevel WishlistPermissionLevel = "guest"
)

var permissionLevelsActions = map[WishlistPermissionLevel]map[WishlistAction]struct{}{
	OwnerWishlistPermissionLevel: {
		ReadWishlistAction:   struct{}{},
		ModifyWishlistAction: struct{}{},
	},
	GuestWishlistPersmissionLevel: {
		ReadWishlistAction:        struct{}{},
		ReserveWishWishlistAction: struct{}{},
	},
}

var (
	ErrWishlistPermissionNotExist = errors.New("wishlist permission does not exist")
)

type WishlistPermission struct {
	ID         int64                   `db:"id"`
	UserID     string                  `db:"user_id"`
	WishlistID string                  `db:"wishlist_id"`
	Level      WishlistPermissionLevel `db:"permission_level"`
}

func NewWishlistPersmission(userID, wishlistID string, level WishlistPermissionLevel) *WishlistPermission {
	return &WishlistPermission{
		UserID:     userID,
		WishlistID: wishlistID,
		Level:      level,
	}
}

func (p *WishlistPermission) Can(action WishlistAction) bool {
	actions, ok := permissionLevelsActions[p.Level]
	if !ok {
		return false
	}
	_, ok = actions[action]
	return ok
}

type WishlistsPermissions []*WishlistPermission

func (p WishlistsPermissions) GetWishlitsIdsForAction(action WishlistAction) []string {
	wishlitsIds := make([]string, 0)
	for _, permission := range p {
		if permission.Can(action) {
			wishlitsIds = append(wishlitsIds, permission.WishlistID)
		}
	}
	return wishlitsIds
}

func (p *WishlistPermission) CanGrantPermission(level WishlistPermissionLevel) bool {
	actions, ok := permissionLevelsActions[p.Level]
	if !ok {
		return false
	}
	_, ok = actions[ModifyWishlistAction]
	if !ok {
		return false
	}
	switch level {
	case GuestWishlistPersmissionLevel:
		return true
	default:
		return false
	}
}

func (p *WishlistPermission) CanRevokePermission(revokingLevel WishlistPermissionLevel) bool {
	actions, ok := permissionLevelsActions[p.Level]
	if !ok {
		return false
	}
	_, ok = actions[ModifyWishlistAction]
	if !ok {
		return false
	}
	switch revokingLevel {
	case GuestWishlistPersmissionLevel:
		return true
	default:
		return false
	}
}
