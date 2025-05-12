package cache

import (
	"time"
)

type CacheOptions struct {
	KeyPrefix         string
	DefaultExpiration time.Duration
}

func DefaultCacheOptions() *CacheOptions {
	return &CacheOptions{
		KeyPrefix:         "app:cache:",
		DefaultExpiration: 5 * time.Minute,
	}
}

func (o *CacheOptions) Validate() error {
	return nil
}

type Option func(*CacheOptions)

func WithKeyPrefix(prefix string) Option {
	return func(o *CacheOptions) {
		o.KeyPrefix = prefix
	}
}

func WithDefaultExpiration(duration time.Duration) Option {
	return func(o *CacheOptions) {
		o.DefaultExpiration = duration
	}
}

func ApplyOptions(o *CacheOptions, opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}
