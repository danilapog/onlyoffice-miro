package deployments

import (
	"embed"
	"fmt"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migrator struct {
	migrate *migrate.Migrate
}

func NewMigrator(pool *pgxpool.Pool, config *config.DataSourceConfig) (*Migrator, error) {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create migrations source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, config.DatasourceURL())
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{
		migrate: m,
	}, nil
}

func (m *Migrator) Up() error {
	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (m *Migrator) Down() error {
	if err := m.migrate.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}

func (m *Migrator) Close() error {
	if _, err := m.migrate.Close(); err != nil {
		return fmt.Errorf("failed to close migrator: %w", err)
	}

	return nil
}
