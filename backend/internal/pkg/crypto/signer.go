package crypto

import "github.com/golang-jwt/jwt/v5"

type Signer interface {
	Validate(tokenString string, secret []byte) (*jwt.Token, error)
	ValidateTarget(tokenString string, secret []byte, target any) error
	Create(claims jwt.Claims, secret []byte) (string, error)
}
