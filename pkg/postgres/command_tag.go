package postgres

import (
	"github.com/jackc/pgx/v5/pgconn"
)

type CommandTag struct {
	*pgconn.CommandTag
}
