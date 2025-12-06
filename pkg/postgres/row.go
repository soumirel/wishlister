package postgres

import "github.com/jackc/pgx/v5"

type Row interface {
	pgx.Row
}

type Rows interface {
	pgx.Rows
}

type row struct {
	Row
}

type rows struct {
	Rows
}
