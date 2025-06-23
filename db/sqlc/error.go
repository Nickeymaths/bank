package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrNoRows = pgx.ErrNoRows
var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

var ErrForeignKeyViolation = &pgconn.PgError{
	Code: ForeignKeyViolation,
}
