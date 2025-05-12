package settings

import (
	"net/url"
	"strings"
)

type SaveOptions struct {
	Address string
	Header  string
	Secret  string
	Demo    bool
}

func (o *SaveOptions) Validate() error {
	if !o.Demo && o.Address == "" {
		return ErrAddressRequired
	}

	if !o.Demo && o.Secret == "" {
		return ErrSecretRequired
	}

	if !o.Demo && o.Header == "" {
		return ErrHeaderRequired
	}

	if o.Address != "" {
		u, err := url.Parse(o.Address)
		if err != nil {
			return ErrInvalidURL
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return ErrInvalidProtocol
		}

		if strings.HasSuffix(o.Address, "/") {
			return ErrTrailingSlash
		}
	}

	if len(o.Header) > 255 {
		return ErrHeaderTooLong
	}

	if len(o.Secret) > 255 {
		return ErrSecretTooLong
	}

	return nil
}

type Option func(*SaveOptions)

func WithAddress(val string) Option {
	return func(o *SaveOptions) {
		o.Address = val
	}
}

func WithHeader(val string) Option {
	return func(o *SaveOptions) {
		o.Header = val
	}
}

func WithSecret(val string) Option {
	return func(o *SaveOptions) {
		o.Secret = val
	}
}

func WithDemo(val bool) Option {
	return func(o *SaveOptions) {
		o.Demo = val
	}
}
