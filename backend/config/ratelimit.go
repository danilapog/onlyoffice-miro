package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type RateLimitConfig struct {
	Rate      int           `yaml:"rate" env:"RATE_LIMIT_RATE" validate:"gt=0"`
	Window    time.Duration `yaml:"window" env:"RATE_LIMIT_WINDOW" validate:"gt=0"`
	SkipPaths []string      `yaml:"skip_paths" env:"RATE_LIMIT_SKIP_PATHS"`
}

func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		Rate:      100,
		Window:    time.Minute,
		SkipPaths: []string{"/health", "/metrics"},
	}
}

func (c *RateLimitConfig) loadEnv() error {
	if rate := os.Getenv("RATE_LIMIT_RATE"); rate != "" {
		if r, err := strconv.Atoi(rate); err == nil {
			c.Rate = r
		}
	}

	if window := os.Getenv("RATE_LIMIT_WINDOW"); window != "" {
		if w, err := time.ParseDuration(window); err == nil {
			c.Window = w
		}
	}

	if skipPaths := os.Getenv("RATE_LIMIT_SKIP_PATHS"); skipPaths != "" {
		c.SkipPaths = strings.Split(skipPaths, ",")
	}

	return nil
}

func (c *RateLimitConfig) Validate() error {
	validate := validator.New()

	if err := validate.Struct(c); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				switch e.Field() {
				case "Rate":
					return fmt.Errorf("rate limit rate must be greater than 0")
				case "Window":
					return fmt.Errorf("rate limit window must be greater than 0")
				default:
					return fmt.Errorf("validation error on field %s: %s", e.Field(), e.Tag())
				}
			}
		}

		return err
	}

	return nil
}
