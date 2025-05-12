package oauth

import "errors"

var (
	ErrTokenMissing = errors.New("token is missing")
	ErrTokenExpired = errors.New("token has expired")
)
