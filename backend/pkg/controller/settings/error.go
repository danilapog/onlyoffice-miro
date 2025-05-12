package settings

import "errors"

var (
	ErrBoardIdRequired       = errors.New("board_id is required")
	ErrInvalidRequestBody    = errors.New("missing or invalid request body")
	ErrMissingBoardParameter = errors.New("missing board id parameter")
	ErrMissingOpenIdToken    = errors.New("oid token is missing")
)
