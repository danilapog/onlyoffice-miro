package oauth

import "errors"

var (
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenMissing = errors.New("token is missing")
)
