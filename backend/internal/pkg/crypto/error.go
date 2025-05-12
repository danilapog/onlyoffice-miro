package crypto

import "errors"

var (
	ErrCipherTextEmpty    = errors.New("cipher text is empty")
	ErrCipherTextTooShort = errors.New("cipher text too short")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrTokenMapping       = errors.New("failed to map token")
)
