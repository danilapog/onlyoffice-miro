package service

import "context"

type Fields map[string]any

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Fields)
	Info(ctx context.Context, msg string, fields ...Fields)
	Warn(ctx context.Context, msg string, fields ...Fields)
	Error(ctx context.Context, msg string, fields ...Fields)
	Fatal(ctx context.Context, msg string, fields ...Fields)
	WithFields(fields Fields) Logger
	WithContext(ctx context.Context) Logger
}
