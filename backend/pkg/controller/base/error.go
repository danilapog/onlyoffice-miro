package base

import "errors"

var (
	ErrMissingAuthentication = errors.New("authentication not found")
	ErrMissingOpenIdToken    = errors.New("oid token is missing")
	ErrSettingsNotConfigured = errors.New("settings are not configured")
)
