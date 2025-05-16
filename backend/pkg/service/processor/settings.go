package processor

import (
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	pgx "github.com/jackc/pgx/v5"
)

const (
	settingsSelectQuery = `SELECT s.address, s.header, s.secret, 
	d.enabled, d.started
	FROM settings s
	LEFT JOIN demos d ON s.team_id = d.team_id
	WHERE s.team_id = $1 AND s.board_id = $2;`

	settingsUpdateQuery = `UPDATE settings
SET address = $3,
    header = $4,
    secret = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE team_id = $1 AND board_id = $2;`

	settingsDeleteQuery = `DELETE FROM settings
WHERE team_id = $1 AND board_id = $2;`
)

type settingsProcessor struct{}

func settingsScanner(row pgx.Row) (*component.Settings, error) {
	result := &component.Settings{}
	var enabled *bool
	var started *time.Time

	if err := row.Scan(
		&result.Address,
		&result.Header,
		&result.Secret,
		&enabled,
		&started,
	); err != nil {
		return nil, err
	}

	if enabled != nil && started != nil {
		result.Demo = component.Demo{
			Enabled: *enabled,
			Started: started,
		}
	}

	return result, nil
}

func NewSettingsProcessor() service.StorageProcessor[core.SettingsCompositeKey, component.Settings, pgx.Row] {
	return &settingsProcessor{}
}

func (s settingsProcessor) TableName() string {
	return "settings"
}

func (s settingsProcessor) BuildSelectQuery(id core.SettingsCompositeKey) (string, []any, func(pgx.Row) (component.Settings, error)) {
	return settingsSelectQuery, []any{id.TeamID, id.BoardID}, func(row pgx.Row) (component.Settings, error) {
		settings, err := settingsScanner(row)
		if err != nil {
			return component.Settings{}, err
		}

		if settings.Demo != (component.Demo{}) {
			settings.Demo.TeamID = id.TeamID
		}

		return *settings, nil
	}
}

func (s settingsProcessor) BuildInsertQuery(id core.SettingsCompositeKey, settings component.Settings) (string, []any) {
	if settings.Demo.Enabled {
		var started *time.Time
		if !settings.Demo.Started.IsZero() {
			started = settings.Demo.Started
		}

		return `
            WITH settings_update AS (
                INSERT INTO settings (team_id, board_id, address, header, secret)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT (team_id, board_id) DO UPDATE
                SET address = EXCLUDED.address,
                    header = EXCLUDED.header,
                    secret = EXCLUDED.secret,
                    updated_at = CURRENT_TIMESTAMP
                RETURNING team_id
            )
            INSERT INTO demos (team_id, enabled, started)
            VALUES ($1, $6, $7)
            ON CONFLICT (team_id) DO NOTHING
            RETURNING team_id
        `, []any{
				id.TeamID,
				id.BoardID,
				settings.Address,
				settings.Header,
				settings.Secret,
				settings.Demo.Enabled,
				started,
			}
	}

	return `
        INSERT INTO settings (team_id, board_id, address, header, secret)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (team_id, board_id) DO UPDATE
        SET address = EXCLUDED.address,
            header = EXCLUDED.header,
            secret = EXCLUDED.secret,
            updated_at = CURRENT_TIMESTAMP
        RETURNING team_id
    `, []any{
			id.TeamID,
			id.BoardID,
			settings.Address,
			settings.Header,
			settings.Secret,
		}
}

func (s settingsProcessor) BuildUpdateQuery(id core.SettingsCompositeKey, settings component.Settings) (string, []any) {
	return settingsUpdateQuery, []any{
		id.TeamID,
		id.BoardID,
		settings.Address,
		settings.Header,
		settings.Secret,
	}
}

func (s settingsProcessor) BuildDeleteQuery(id core.SettingsCompositeKey) (string, []any) {
	return settingsDeleteQuery, []any{id.TeamID, id.BoardID}
}
