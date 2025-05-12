package logger

import (
	"context"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
)

type NoopLogger struct{}

func NewNoopLogger() service.Logger {
	return &NoopLogger{}
}

func (l *NoopLogger) Debug(_ context.Context, _ string, _ ...service.Fields) {}
func (l *NoopLogger) Info(_ context.Context, _ string, _ ...service.Fields)  {}
func (l *NoopLogger) Warn(_ context.Context, _ string, _ ...service.Fields)  {}
func (l *NoopLogger) Error(_ context.Context, _ string, _ ...service.Fields) {}
func (l *NoopLogger) Fatal(_ context.Context, _ string, _ ...service.Fields) {}

func (l *NoopLogger) WithFields(_ service.Fields) service.Logger {
	return l
}

func (l *NoopLogger) WithContext(_ context.Context) service.Logger {
	return l
}
