package settings

import "errors"

var (
	ErrAddressRequired = errors.New("address is required in non-demo mode")
	ErrSecretRequired  = errors.New("secret is required in demo mode")
	ErrHeaderRequired  = errors.New("header is required in demo mode")
	ErrInvalidURL      = errors.New("address must be a valid URL")
	ErrInvalidProtocol = errors.New("address must use http or https protocol")
	ErrTrailingSlash   = errors.New("address must not have a trailing slash")
	ErrHeaderTooLong   = errors.New("header must be at most 255 characters")
	ErrSecretTooLong   = errors.New("secret must be at most 255 characters")
)
