package config

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)

type CookieConfig struct {
	Name     string `yaml:"name" env:"COOKIE_NAME" validate:"required"`
	Path     string `yaml:"path" env:"COOKIE_PATH" validate:"required"`
	MaxAge   int    `yaml:"max_age" env:"COOKIE_MAX_AGE" validate:"gt=0"`
	Secure   bool   `yaml:"secure" env:"COOKIE_SECURE"`
	HttpOnly bool   `yaml:"http_only" env:"COOKIE_HTTP_ONLY"`
	SameSite string `yaml:"same_site" env:"COOKIE_SAME_SITE" validate:"oneof=None Lax Strict"`
}

func DefaultCookieConfig() *CookieConfig {
	return &CookieConfig{
		Name:     "asc_miro_token",
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		Secure:   true,
		HttpOnly: true,
		SameSite: "None",
	}
}

func (c *CookieConfig) loadEnv() error {
	if name := os.Getenv("COOKIE_NAME"); name != "" {
		c.Name = name
	}

	if path := os.Getenv("COOKIE_PATH"); path != "" {
		c.Path = path
	}

	if maxAge := os.Getenv("COOKIE_MAX_AGE"); maxAge != "" {
		if duration, err := time.ParseDuration(maxAge); err != nil {
			return fmt.Errorf("invalid max_age duration: %w", err)
		} else {
			c.MaxAge = int(duration.Seconds())
		}
	}

	if secure := os.Getenv("COOKIE_SECURE"); secure != "" {
		c.Secure = secure == "true"
	}

	if httpOnly := os.Getenv("COOKIE_HTTP_ONLY"); httpOnly != "" {
		c.HttpOnly = httpOnly == "true"
	}

	if sameSite := os.Getenv("COOKIE_SAME_SITE"); sameSite != "" {
		c.SameSite = sameSite
	}

	return nil
}

func (c *CookieConfig) Validate() error {
	validate := validator.New()

	if err := validate.Struct(c); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				switch e.Field() {
				case "Name":
					return fmt.Errorf("cookie name is required")
				case "Path":
					return fmt.Errorf("cookie path is required")
				case "MaxAge":
					return fmt.Errorf("cookie max_age must be positive")
				case "SameSite":
					return fmt.Errorf("invalid same_site value: %s", c.SameSite)
				default:
					return fmt.Errorf("validation error on field %s: %s", e.Field(), e.Tag())
				}
			}
		}

		return err
	}

	return nil
}

func (c *CookieConfig) GetSameSite() http.SameSite {
	switch c.SameSite {
	case "None":
		return http.SameSiteNoneMode
	case "Lax":
		return http.SameSiteLaxMode
	case "Strict":
		return http.SameSiteStrictMode
	default:
		return http.SameSiteNoneMode
	}
}
