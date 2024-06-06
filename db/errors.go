package db

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
