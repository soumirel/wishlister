package useridentity

type GetUserIdByExternalIdentityCommand struct {
	ExternalID       string
	IdentityProvider string
}

type LinkUserWithExternalIdentityCommand struct {
	UserID           string
	ExternalID       string
	IdentityProvider string
}

type CreateUserFromExternalIdentityCommand struct {
	ExternalID       string
	IdentityProvider string
}
