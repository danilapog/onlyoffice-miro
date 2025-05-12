package logger

import "errors"

var (
	ErrInvalidLogLevel = errors.New("invalid log level")
	ErrNilContext      = errors.New("nil context provided")
)
