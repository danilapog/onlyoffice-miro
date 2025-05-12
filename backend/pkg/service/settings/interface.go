package settings

import (
	"context"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
)

type SettingsService interface {
	Save(ctx context.Context, teamID, boardID string, opts ...Option) error
	Find(ctx context.Context, teamID, boardID string) (component.Settings, error)
}
