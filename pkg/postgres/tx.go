package postgres

import (
	"github.com/jackc/pgx/v5"
)

type Tx interface {
	pgx.Tx
}

type tx struct {
	Tx
}
