package auth

import (
	"context"
	"errors"
)

type CtxKey int

const (
	_ CtxKey = iota
	AuthCtxKey
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Auth struct {
	UserID string
}

func NewCtx(ctx context.Context, a Auth) context.Context {
	return context.WithValue(ctx, AuthCtxKey, a)
}

func FromCtx(ctx context.Context) (Auth, bool) {
	aauth := ctx.Value(AuthCtxKey)
	auth, ok := aauth.(Auth)
	return auth, ok
}

func FromCtxOrEmpty(ctx context.Context) Auth {
	au, ok := FromCtx(ctx)
	if !ok {
		return Auth{}
	}
	return au
}
