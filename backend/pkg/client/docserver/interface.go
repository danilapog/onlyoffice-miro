package docserver

import "context"

type Client interface {
	GetServerVersion(ctx context.Context, base string, opts ...Option) (*ServerVersionResponse, error)
}
