package callback

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/common"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/oauth"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/settings"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

type callbackController struct {
	config          *config.Config
	miroClient      miro.Client
	oauthService    oauth.OAuthService[miro.AuthenticationResponse]
	settingsService settings.SettingsService
	jwtService      crypto.Signer
	logger          service.Logger
}

func NewCallbackController(
	config *config.Config,
	miroClient miro.Client,
	oauthService oauth.OAuthService[miro.AuthenticationResponse],
	settingsService settings.SettingsService,
	jwtService crypto.Signer,
	logger service.Logger,
) common.Handler {
	controller := &callbackController{
		config:          config,
		oauthService:    oauthService,
		miroClient:      miroClient,
		settingsService: settingsService,
		jwtService:      jwtService,
		logger:          logger,
	}

	return common.NewHandler(map[common.HTTPMethod]echo.HandlerFunc{
		common.MethodPost: controller.handlePost,
	})
}

func (c *callbackController) logErrorAndRespond(ctx echo.Context, statusCode int, logMessage string, err error) error {
	c.logger.Error(ctx.Request().Context(), logMessage, service.Fields{"error": err.Error()})
	return ctx.JSON(statusCode, common.ErrorResponse{Error: callbackErrorCodeFailure})
}

func (c *callbackController) extractParams(ctx echo.Context) (callbackQueryParams, error) {
	params := callbackQueryParams{
		UID: ctx.QueryParam("uid"),
		TID: ctx.QueryParam("tid"),
		BID: ctx.QueryParam("bid"),
		FID: ctx.QueryParam("fid"),
	}

	if params.UID == "" || params.BID == "" || params.TID == "" || params.FID == "" {
		return params, fmt.Errorf("missing required query parameters")
	}

	return params, nil
}

func (c *callbackController) fetchAuthenticationAndSettings(
	ctx context.Context,
	params callbackQueryParams,
) (component.Authentication, component.Settings, error) {
	var settings component.Settings
	var auth component.Authentication

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		settings, err = c.settingsService.Find(ctx, params.TID, params.BID)
		if err != nil {
			return fmt.Errorf("failed to fetch settings: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		auth, err = c.oauthService.Find(ctx, params.TID, params.UID)
		if err != nil {
			return fmt.Errorf("failed to fetch auth: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return auth, settings, err
	}

	return auth, settings, nil
}

func (c *callbackController) getSecretFromSettings(settings component.Settings) string {
	if settings.Demo.Enabled &&
		settings.Address == "" &&
		settings.Demo.Started.Add(time.Duration(c.config.DemoServer.Days)*24*time.Hour).After(time.Now()) {
		return c.config.DemoServer.Secret
	}
	return settings.Secret
}

func (c *callbackController) handlePost(ctx echo.Context) error {
	tctx, cancel := context.WithTimeout(ctx.Request().Context(), saveFileRequestTimeout)
	defer cancel()

	params, err := c.extractParams(ctx)
	if err != nil {
		return c.logErrorAndRespond(ctx, http.StatusBadRequest, "failed to extract query parameters", err)
	}

	var body callbackRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&body); err != nil {
		return c.logErrorAndRespond(ctx, http.StatusBadRequest, "failed to decode request body", err)
	}

	if err := body.Validate(); err != nil {
		return c.logErrorAndRespond(ctx, http.StatusBadRequest, "failed to validate request body", err)
	}

	if body.Status == 2 {
		if body.Token == "" {
			c.logger.Debug(ctx.Request().Context(), "failed to extract token from request body", nil)
			return ctx.JSON(http.StatusBadRequest, common.ErrorResponse{Error: callbackErrorCodeFailure})
		}

		auth, settings, err := c.fetchAuthenticationAndSettings(tctx, params)
		if err != nil {
			return c.logErrorAndRespond(ctx, http.StatusBadRequest, "failed to extract authentication and settings", err)
		}

		secret := c.getSecretFromSettings(settings)
		if err = c.jwtService.ValidateTarget(body.Token, []byte(secret), &body); err != nil {
			return c.logErrorAndRespond(ctx, http.StatusUnauthorized, "failed to validate and map token", err)
		}

		if _, err := c.miroClient.UploadFile(tctx, miro.UploadFileRequest{
			BoardID: params.BID,
			ItemID:  params.FID,
			FileURL: body.Url,
			Token:   auth.AccessToken,
		}); err != nil {
			c.logger.Error(ctx.Request().Context(), "failed to upload file",
				service.Fields{
					"error":    err.Error(),
					"board_id": params.BID,
					"file_id":  params.FID,
				},
			)
			return ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: callbackErrorCodeFailure})
		}

		c.logger.Info(ctx.Request().Context(), "file uploaded successfully",
			service.Fields{
				"board_id": params.BID,
				"file_id":  params.FID,
			},
		)
	}

	return ctx.JSON(http.StatusOK, common.ErrorResponse{Error: callbackErrorCodeSuccess})
}
