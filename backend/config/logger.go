package config

import (
	"io"
	"os"
)

type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
	Fatal Level = "fatal"
)

type LoggerConfig struct {
	ServiceName string    `yaml:"service_name"`
	Environment string    `yaml:"environment"`
	Level       string    `yaml:"level"`
	PrettyPrint bool      `yaml:"pretty_print"`
	LoggerType  string    `yaml:"logger_type"`
	Output      io.Writer `yaml:"-"`
}

func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		ServiceName: "onlyoffice-miro-service",
		Environment: "development",
		Level:       "info",
		PrettyPrint: false,
		LoggerType:  "zap",
		Output:      os.Stdout,
	}
}

func (c *LoggerConfig) loadEnv() error {
	if val, exists := os.LookupEnv("LOGGER_SERVICE_NAME"); exists {
		c.ServiceName = val
	}

	if val, exists := os.LookupEnv("LOGGER_ENVIRONMENT"); exists {
		c.Environment = val
	}

	if val, exists := os.LookupEnv("LOGGER_LEVEL"); exists {
		c.Level = val
	}

	if val, exists := os.LookupEnv("LOGGER_PRETTY_PRINT"); exists {
		c.PrettyPrint = val == "true"
	}

	if val, exists := os.LookupEnv("LOGGER_TYPE"); exists {
		c.LoggerType = val
	}

	return nil
}

func (c *LoggerConfig) Validate() error {
	return nil
}

func (c *LoggerConfig) ToLogLevel() Level {
	switch c.Level {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	case "fatal":
		return Fatal
	default:
		return Info
	}
}
