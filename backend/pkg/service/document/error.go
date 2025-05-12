package document

import "errors"

var (
	ErrCallbackNotFound  = errors.New("callback handler does not exist")
	ErrUnsupportedFormat = errors.New("unsupported file format")
)
