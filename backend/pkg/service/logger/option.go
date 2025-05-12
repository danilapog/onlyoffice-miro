package logger

import (
	"io"
	"os"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
)

type Option func(*config.LoggerConfig)

func WithServiceName(name string) Option {
	return func(c *config.LoggerConfig) {
		c.ServiceName = name
	}
}

func WithEnvironment(env string) Option {
	return func(c *config.LoggerConfig) {
		c.Environment = env
	}
}

func WithLevel(level config.Level) Option {
	return func(c *config.LoggerConfig) {
		c.Level = string(level)
	}
}

func WithPrettyPrint(enabled bool) Option {
	return func(c *config.LoggerConfig) {
		c.PrettyPrint = enabled
	}
}

func WithOutput(w io.Writer) Option {
	return func(c *config.LoggerConfig) {
		c.Output = w
	}
}

func NewLoggerWithOptions(loggerType string, opts ...Option) service.Logger {
	config := config.DefaultLoggerConfig()

	for _, opt := range opts {
		opt(config)
	}

	if config.Output == nil {
		config.Output = os.Stdout
	}

	switch loggerType {
	case "zap":
		return NewZapLogger(config)
	case "noop":
		return NewNoopLogger()
	default:
		return NewZapLogger(config)
	}
}
