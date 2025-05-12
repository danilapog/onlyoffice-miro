package initializer

import (
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/oauth"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/common"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/controller"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/document"
	oauthService "github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/oauth"
	settingsService "github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/settings"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// App is the main application container that holds all dependencies.
type App struct {
	Config      *config.Config
	Clients     *Clients
	Controllers *Controllers
	Database    *Database
	Echo        *echo.Echo
	Router      *Router
	Services    *Services
}

// SetupRoutes calls the router's route setup method.
func (a *App) SetupRoutes(logger service.Logger) {
	// Route setup is now implemented in router.go
	a.Router.SetupRoutes(a.Controllers, logger)
}

// Database holds all database-related components and repositories.
type Database struct {
	AuthStorage     service.Storage[core.AuthCompositeKey, component.Authentication]
	Pool            *pgxpool.Pool
	SettingsStorage service.Storage[core.SettingsCompositeKey, component.Settings]
}

// Clients contains all external API client instances.
type Clients struct {
	OAuthClient oauth.OAuthClient[miro.AuthenticationResponse]
	MiroClient  miro.Client
}

// Services contains all application service instances.
type Services struct {
	AuthService     oauthService.OAuthService[miro.AuthenticationResponse]
	Builder         document.BuilderService
	FormatManager   document.FormatManager
	JwtService      crypto.Signer
	Renderer        *controller.TemplateRenderer
	SettingsService settingsService.SettingsService
}

// Controllers contains all HTTP request handlers.
type Controllers struct {
	Auth           common.Handler
	Callback       common.Handler
	Editor         common.Handler
	FileConversion common.Handler
	FileManagement common.Handler
	Settings       common.Handler
}

// Router provides access to the Echo instance and configuration.
type Router struct {
	Config   *config.Config
	Echo     *echo.Echo
	Services *Services
}
