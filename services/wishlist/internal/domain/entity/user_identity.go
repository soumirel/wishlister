package entity

import (
	"errors"

	"github.com/gofrs/uuid/v5"
)

type IdentityProvider string

const (
	telegramIdentityProvider = "telegram"
)

var (
	ErrUserIdentityDoesNotExist     = errors.New("user identity does not exist")
	ErrExternalUserIdentityConflict = errors.New("external user identity conflict")
)

type ExternalIdentity struct {
	ExternalID string
	Provider   IdentityProvider
}

type UserIdentity struct {
	ID               string
	UserID           string
	ExternalIdentity ExternalIdentity
}

func NewUserIdentity(userID string, externalIdentity ExternalIdentity) *UserIdentity {
	return &UserIdentity{
		ID:               uuid.Must(uuid.NewV4()).String(),
		UserID:           userID,
		ExternalIdentity: externalIdentity,
	}
}

func NewExternalIdentity(externalID string, identityProvider IdentityProvider) (ExternalIdentity, error) {
	return ExternalIdentity{
		ExternalID: externalID,
		Provider:   identityProvider,
	}, nil
}
