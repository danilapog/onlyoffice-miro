package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/initializer"
	fx "go.uber.org/fx"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = filepath.Join("config.yaml")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}

	if err := cfg.Validate(); err != nil {
		slog.Error("invalid config", "error", err)
		return
	}

	app := fx.New(
		fx.NopLogger,
		fx.Provide(func() *config.Config { return cfg }),
		fx.Provide(func() *config.LoggerConfig { return cfg.Logger }),
		initializer.Module,
		fx.Invoke(func(app *initializer.App) {
			if err := app.Echo.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("failed to start server", "error", err)
			}
		}),
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		slog.Error("failed to start application", "error", err)
		return
	}

	<-app.Done()
}
