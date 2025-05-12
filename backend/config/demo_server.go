package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type DemoServerConfig struct {
	Address string `yaml:"address" env:"DEMO_SERVER_ADDRESS" validate:"required,http_address"`
	Header  string `yaml:"header" env:"DEMO_SERVER_HEADER" validate:"required"`
	Secret  string `yaml:"secret" env:"DEMO_SERVER_SECRET" validate:"required"`
	Days    int    `yaml:"days" env:"DEMO_SERVER_DAYS" validate:"gt=0"`
}

func DefaultDemoServerConfig() *DemoServerConfig {
	return &DemoServerConfig{
		Address: "http://localhost:8080",
		Header:  "AuthorizationJwt",
		Secret:  "secret",
		Days:    30,
	}
}

func (c *DemoServerConfig) loadEnv() error {
	if address := os.Getenv("DEMO_SERVER_ADDRESS"); address != "" {
		c.Address = address
	}

	if header := os.Getenv("DEMO_SERVER_HEADER"); header != "" {
		c.Header = header
	}

	if secret := os.Getenv("DEMO_SERVER_SECRET"); secret != "" {
		c.Secret = secret
	}

	if days := os.Getenv("DEMO_SERVER_DAYS"); days != "" {
		if daysInt, err := strconv.Atoi(days); err == nil {
			c.Days = daysInt
		}
	}

	return nil
}

func (c *DemoServerConfig) Validate() error {
	validate := validator.New()

	validate.RegisterValidation("http_address", func(fl validator.FieldLevel) bool {
		address := fl.Field().String()
		if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
			return false
		}

		if strings.HasSuffix(address, "/") {
			return false
		}

		return true
	})

	if err := validate.Struct(c); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				switch e.Field() {
				case "Address":
					if e.Tag() == "http_address" {
						return fmt.Errorf("address must be an HTTP address without trailing slash")
					}
					return fmt.Errorf("address is required")
				case "Header":
					return fmt.Errorf("header is required")
				case "Secret":
					return fmt.Errorf("secret is required")
				case "Days":
					return fmt.Errorf("days must be greater than 0")
				default:
					return fmt.Errorf("validation error on field %s: %s", e.Field(), e.Tag())
				}
			}
		}

		return err
	}

	return nil
}
