package file

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"path"
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
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type fileConversionController struct {
	*base.BaseController
}

func NewFileConversionController(
	config *config.Config,
	miroClient miro.Client,
	builderService document.BuilderService,
	oauthService oauthService.OAuthService[miro.AuthenticationResponse],
	settingsService settings.SettingsService,
	jwtService crypto.Signer,
	logger service.Logger,
) common.Handler {
	controller := &fileConversionController{
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
		common.MethodGet: controller.handleGet,
	})
}

func (c *fileConversionController) handleGet(ctx echo.Context) error {
	return c.ExecuteWithTimeout(ctx, 2*time.Second, func(tctx context.Context) error {
		boardAuth, err := PrepareRequest(ctx, tctx, c.BaseController)
		if err != nil {
			return err
		}

		if boardAuth == nil {
			return c.HandleError(ctx, fmt.Errorf("board authentication is nil"), http.StatusInternalServerError, "invalid board authentication")
		}

		if boardAuth.Authentication == nil {
			return c.HandleError(ctx, fmt.Errorf("authentication data is nil"), http.StatusInternalServerError, "missing authentication data")
		}

		if fid, ferr := c.GetQueryParam(ctx, "fid"); ferr == nil {
			file, err := GetFileInfo(ctx, tctx, c.BaseController, boardAuth.BoardID, fid, boardAuth.Authentication.AccessToken)
			if err != nil {
				return err
			}

			location, err := c.MiroClient.GetFilePublicURL(tctx, miro.GetFilePublicURLRequest{
				URL:   file.Data.DocumentURL,
				Token: boardAuth.Authentication.AccessToken,
			})

			if err != nil {
				return c.HandleError(ctx, err, http.StatusInternalServerError, "failed to fetch miro file")
			}

			token, err := c.ExtractUserToken(ctx)
			if err != nil {
				return c.HandleError(ctx, err, http.StatusForbidden, "failed to extract token")
			}

			settings, err := c.SettingsService.Find(tctx, token.Team, boardAuth.BoardID)
			if err != nil {
				return c.HandleError(ctx, err, http.StatusBadRequest, "failed to fetch settings")
			}

			address := settings.Address
			secret := settings.Secret
			if settings.Demo.Enabled && settings.Demo.Started.Add(time.Duration(c.Config.DemoServer.Days)*24*time.Hour).After(time.Now()) && (address == "" || secret == "") {
				address = c.Config.DemoServer.Address
				secret = c.Config.DemoServer.Secret
			}

			fileExt := path.Ext(file.Data.Title)
			convReq := convertClaims{
				Async:      false,
				FileType:   string(toDocumentType(fileExt)),
				Key:        fmt.Sprintf("%x", md5.Sum([]byte(file.Data.DocumentURL))),
				OutputType: "pdf",
				Title:      file.Data.Title,
				URL:        location.URL,
				RegisteredClaims: jwt.RegisteredClaims{
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
				},
			}

			jwtToken, err := c.JwtService.Create(convReq, []byte(secret))
			if err != nil {
				return c.HandleError(ctx, err, http.StatusBadRequest, "failed to create token")
			}

			return ctx.JSON(200, convertResponse{
				URL:   address,
				Token: jwtToken,
			})
		}

		return ctx.JSON(200, nil)
	})
}
