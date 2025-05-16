package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Database   *DataSourceConfig `yaml:"database"`
	Miro       *MiroConfig       `yaml:"miro"`
	OAuth      *OAuthConfig      `yaml:"oauth"`
	Server     *ServerConfig     `yaml:"server"`
	Redis      *RedisConfig      `yaml:"redis"`
	RateLimit  *RateLimitConfig  `yaml:"rate_limit"`
	Cookie     *CookieConfig     `yaml:"cookie"`
	CORS       *CORSConfig       `yaml:"cors"`
	DemoServer *DemoServerConfig `yaml:"demo_server"`
	Logger     *LoggerConfig     `yaml:"logger"`
}

func DefaultConfig() *Config {
	return &Config{
		Database:   DefaultDataSourceConfig(),
		Miro:       DefaultMiroConfig(),
		OAuth:      DefaultOAuthConfig(),
		Server:     DefaultServerConfig(),
		Redis:      DefaultRedisConfig(),
		RateLimit:  DefaultRateLimitConfig(),
		Cookie:     DefaultCookieConfig(),
		CORS:       DefaultCORSConfig(),
		DemoServer: DefaultDemoServerConfig(),
		Logger:     DefaultLoggerConfig(),
	}
}

func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig()
	if path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			if err := yaml.Unmarshal(data, &config); err != nil {
				return config, fmt.Errorf("failed to parse YAML config: %w", err)
			}
		}
	}

	if err := config.Database.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load database environment variables: %w", err)
	}

	if err := config.Miro.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load Miro environment variables: %w", err)
	}

	if err := config.OAuth.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load OAuth environment variables: %w", err)
	}

	if err := config.Server.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load server environment variables: %w", err)
	}

	if err := config.Redis.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load Redis environment variables: %w", err)
	}

	if err := config.RateLimit.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load rate limit environment variables: %w", err)
	}

	if err := config.Cookie.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load cookie environment variables: %w", err)
	}

	if err := config.CORS.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load CORS environment variables: %w", err)
	}

	if err := config.DemoServer.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load demo server environment variables: %w", err)
	}

	if err := config.Logger.loadEnv(); err != nil {
		return config, fmt.Errorf("failed to load logger environment variables: %w", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("invalid database config: %w", err)
	}

	if err := c.Miro.Validate(); err != nil {
		return fmt.Errorf("invalid Miro config: %w", err)
	}

	if err := c.OAuth.Validate(); err != nil {
		return fmt.Errorf("invalid OAuth config: %w", err)
	}

	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("invalid server config: %w", err)
	}

	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("invalid Redis config: %w", err)
	}

	if err := c.RateLimit.Validate(); err != nil {
		return fmt.Errorf("invalid rate limit config: %w", err)
	}

	if err := c.CORS.Validate(); err != nil {
		return fmt.Errorf("invalid CORS config: %w", err)
	}

	if err := c.DemoServer.Validate(); err != nil {
		return fmt.Errorf("invalid demo server config: %w", err)
	}

	if err := c.Logger.Validate(); err != nil {
		return fmt.Errorf("invalid logger config: %w", err)
	}

	return nil
}
