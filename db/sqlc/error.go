package db

import (
	"database/sql"

	"github.com/lib/pq"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = sql.ErrNoRows

var ErrUniqueViolation = &pq.Error{
	Code: UniqueViolation,
}

var ErrForeignKeyViolation = &pq.Error{
	Code: ForeignKeyViolation,
}
