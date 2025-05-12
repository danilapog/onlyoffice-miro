package oauth

import "context"

type OAuthClient[T any] interface {
	Exchange(ctx context.Context, code string) (T, error)
	Refresh(ctx context.Context, refreshToken string) (T, error)
}
