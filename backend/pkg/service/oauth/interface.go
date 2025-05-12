package oauth

import (
	"context"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
)

type OAuthService[T any] interface {
	Save(ctx context.Context, teamID, userID string, token component.Authentication) error
	Find(ctx context.Context, teamID, userID string) (component.Authentication, error)
}
