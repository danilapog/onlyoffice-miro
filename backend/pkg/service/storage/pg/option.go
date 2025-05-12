package pg

import (
	"runtime"
	"time"
)

type Options struct {
	MaxConnections              int
	MinConnections              int
	ConnectionTimeout           time.Duration
	ConnectionIdleTimeout       time.Duration
	ConnectionHealthcheckPeriod time.Duration
	MaxConnLifetime             time.Duration
	MaxRetries                  int
	RetryInterval               time.Duration
}

type Option func(*Options)

func WithMaxConnections(val int) Option {
	return func(o *Options) {
		if val <= 0 || val > runtime.NumCPU() {
			o.MaxConnections = runtime.NumCPU()
			return
		}

		o.MaxConnections = val
	}
}

func WithMinConnections(val int) Option {
	return func(o *Options) {
		if val <= 0 || val > runtime.NumCPU() {
			o.MinConnections = 1
			return
		}

		o.MinConnections = val
	}
}

func WithConnectionTimeout(val time.Duration) Option {
	return func(o *Options) {
		o.ConnectionTimeout = val
	}
}

func WithConnectionIdleTimeout(val time.Duration) Option {
	return func(o *Options) {
		o.ConnectionIdleTimeout = val
	}
}

func WithConnectionHealthcheckPeriod(val time.Duration) Option {
	return func(o *Options) {
		o.ConnectionHealthcheckPeriod = val
	}
}

func WithMaxConnLifetime(val time.Duration) Option {
	return func(o *Options) {
		o.MaxConnLifetime = val
	}
}

func WithMaxRetries(val int) Option {
	return func(o *Options) {
		if val < 0 {
			o.MaxRetries = 0
			return
		}

		o.MaxRetries = val
	}
}

func WithRetryInterval(val time.Duration) Option {
	return func(o *Options) {
		if val < 100*time.Millisecond {
			o.RetryInterval = 100 * time.Millisecond
			return
		}

		o.RetryInterval = val
	}
}
