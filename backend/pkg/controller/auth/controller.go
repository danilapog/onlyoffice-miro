package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/oauth"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/common"
	oauthService "github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/oauth"
	"github.com/labstack/echo/v4"
)

type authController struct {
	installationLocation string
	oauthClient          oauth.OAuthClient[miro.AuthenticationResponse]
	oauthService         oauthService.OAuthService[miro.AuthenticationResponse]
	logger               service.Logger
}

func NewAuthController(
	config *config.Config,
	oauthClient oauth.OAuthClient[miro.AuthenticationResponse],
	oauthService oauthService.OAuthService[miro.AuthenticationResponse],
	logger service.Logger,
) common.Handler {
	controller := &authController{
		installationLocation: fmt.Sprintf(
			"%s?response_type=code&client_id=%s&redirect_uri=%s",
			miroInstallationBase, config.OAuth.ClientID, config.OAuth.RedirectURI,
		),
		oauthClient:  oauthClient,
		oauthService: oauthService,
		logger:       logger,
	}

	return common.NewHandler(map[common.HTTPMethod]echo.HandlerFunc{
		common.MethodGet: controller.handleGet,
	})
}

func (c *authController) extractParams(ctx echo.Context) (authQueryParams, error) {
	params := authQueryParams{
		Code: ctx.QueryParam("code"),
	}

	if params.Code == "" {
		return params, fmt.Errorf("missing authorization code")
	}

	return params, nil
}

func (c *authController) handleError(ctx echo.Context, msg string, err error, args ...any) error {
	fields := service.Fields{"error": err.Error()}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[args[i].(string)] = args[i+1]
		}
	}

	c.logger.Error(ctx.Request().Context(), msg, fields)
	return ctx.Render(http.StatusBadRequest, "exchange", map[string]any{
		"installationLocation": c.installationLocation,
	})
}

func (c *authController) handleGet(ctx echo.Context) error {
	tctx, cancel := context.WithTimeout(ctx.Request().Context(), 3*time.Second)
	defer cancel()

	params, err := c.extractParams(ctx)
	if err != nil {
		return c.handleError(ctx, "failed to extract authorization code", err)
	}

	token, err := c.oauthClient.Exchange(tctx, params.Code)
	if err != nil {
		return c.handleError(ctx, "failed to exchange authorization code", err,
			"code", params.Code)
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(token.ExpiresIn-10)).Unix()
	auth := component.Authentication{
		TokenType:    token.TokenType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    int(expiresAt),
		Scope:        token.Scope,
	}

	if err := c.oauthService.Save(tctx, token.TeamID, token.UserID, auth); err != nil {
		return c.handleError(ctx, "failed to persist authentication token", err,
			"user_id", token.UserID,
			"team_id", token.TeamID)
	}

	c.logger.Info(ctx.Request().Context(), "successfully authenticated user",
		service.Fields{
			"user_id": token.UserID,
			"team_id": token.TeamID,
		},
	)

	return ctx.Redirect(http.StatusPermanentRedirect, miroApplicationBase)
}
