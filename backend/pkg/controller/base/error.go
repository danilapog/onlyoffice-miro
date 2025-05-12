package base

import "errors"

var (
	ErrMissingOpenIdToken    = errors.New("oid token is missing")
	ErrSettingsNotConfigured = errors.New("settings are not configured")
	ErrMissingAuthentication = errors.New("authentication not found")
)
