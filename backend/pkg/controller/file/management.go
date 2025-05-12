package file

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/common"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/controller/base"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/document"
	oauthService "github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/oauth"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/settings"
	"github.com/labstack/echo/v4"
)

type fileManagementController struct {
	*base.BaseController
}

func NewFileManagementController(
	config *config.Config,
	miroClient miro.Client,
	builderService document.BuilderService,
	oauthService oauthService.OAuthService[miro.AuthenticationResponse],
	settingsService settings.SettingsService,
	jwtService crypto.Signer,
	logger service.Logger,
) common.Handler {
	controller := &fileManagementController{
		BaseController: base.NewBaseController(
			config,
			miroClient,
			builderService,
			oauthService,
			settingsService,
			jwtService,
			logger,
		),
	}

	return common.NewHandler(map[common.HTTPMethod]echo.HandlerFunc{
		common.MethodGet:  controller.handleGet,
		common.MethodPost: controller.handlePost,
	})
}

func (c *fileManagementController) handleGet(ctx echo.Context) error {
	return c.ExecuteWithTimeout(ctx, 4*time.Second, func(tctx context.Context) error {
		boardAuth, err := PrepareRequest(ctx, tctx, c.BaseController)
		if err != nil {
			return err
		}

		if fid, ferr := c.GetQueryParam(ctx, "fid"); ferr == nil {
			file, err := GetFileInfo(ctx, tctx, c.BaseController, boardAuth.BoardID, fid, boardAuth.Authentication.AccessToken)
			if err != nil {
				return err
			}
			return ctx.JSON(200, file)
		}

		var cursor string
		if c := ctx.QueryParam("cursor"); c != "" {
			cursor = c
		}

		files, err := GetFilesInfo(ctx, tctx, c.BaseController, boardAuth.BoardID, cursor, boardAuth.Authentication.AccessToken)
		if err != nil {
			return err
		}

		return ctx.JSON(200, files)
	})
}

func (c *fileManagementController) handlePost(ctx echo.Context) error {
	return c.ExecuteWithTimeout(ctx, 4*time.Second, func(tctx context.Context) error {
		var body createBody
		if err := json.NewDecoder(ctx.Request().Body).Decode(&body); err != nil {
			return c.HandleError(ctx, err, http.StatusBadRequest, "failed to decode request body")
		}

		token, err := c.ExtractUserToken(ctx)
		if err != nil {
			return c.HandleError(ctx, err, http.StatusBadRequest, "failed to extract authentication parameters")
		}

		_, auth, err := c.FetchAuthenticationWithSettings(tctx, token.User, token.Team, body.BoardId)
		if err != nil {
			return c.HandleError(ctx, err, http.StatusBadRequest, "failed to fetch required data")
		}

		req := miro.CreateFileRequest{
			BoardID:  body.BoardId,
			Name:     body.FileName,
			Type:     toDocumentType(body.FileType),
			Language: common.ToTemplateLanguage(body.FileLang),
			Token:    auth.AccessToken,
		}

		response, err := CreateFile(ctx, tctx, c.BaseController, req)
		if err != nil {
			return err
		}

		return c.SendJSON(ctx, response)
	})
}
