package docserver

type ClientOptions struct {
	Token  string
	Header string
}

func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Token:  "",
		Header: "",
	}
}

func (o *ClientOptions) Validate() error {
	return nil
}

type Option func(*ClientOptions)

func WithToken(token string) Option {
	return func(o *ClientOptions) {
		o.Token = token
	}
}

func WithHeader(header string) Option {
	return func(o *ClientOptions) {
		o.Header = header
	}
}

func ApplyOptions(o *ClientOptions, opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}
