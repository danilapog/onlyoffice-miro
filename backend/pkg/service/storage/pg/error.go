package pg

import "errors"

var (
	ErrNilPool        = errors.New("received a nil pgx pool")
	ErrNoRowsAffected = errors.New("no rows have been affected")
)
