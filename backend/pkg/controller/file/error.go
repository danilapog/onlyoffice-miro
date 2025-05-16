package file

import "errors"

var (
	ErrInvalidBoardAuthentication = errors.New("invalid board authentication")
	ErrMissingAuthenticationData  = errors.New("missing authentication data")
	ErrFailedToFetchMiroFile      = errors.New("failed to fetch miro file")
	ErrFailedToExtractToken       = errors.New("failed to extract token")
	ErrFailedToFetchSettings      = errors.New("failed to fetch settings")
)
