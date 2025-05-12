package editor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/core/component"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/crypto"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/client/miro"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/common"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/controller/base"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/document"
	oauthService "github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/oauth"
	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/service/settings"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

type editorController struct {
	*base.BaseController
}

func NewEditorController(
	config *config.Config,
	miroClient miro.Client,
	builderService document.BuilderService,
	oauthService oauthService.OAuthService[miro.AuthenticationResponse],
	settingsService settings.SettingsService,
	jwtService crypto.Signer,
	logger service.Logger,
) common.Handler {
	controller := &editorController{
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

func buildCallbackURL(base, fid, uid, tid, bid string) string {
	return fmt.Sprintf("%s?fid=%s&uid=%s&tid=%s&bid=%s", base, fid, uid, tid, bid)
}

func (c *editorController) extractAndValidateParams(ctx echo.Context) (editorRequestParams, error) {
	token, err := c.ExtractUserToken(ctx)
	if err != nil {
		return editorRequestParams{}, err
	}

	bid, err := c.GetQueryParam(ctx, "bid")
	if err != nil {
		return editorRequestParams{}, err
	}

	fid, err := c.GetQueryParam(ctx, "fid")
	if err != nil {
		return editorRequestParams{}, err
	}

	lang, err := c.GetQueryParam(ctx, "lang")
	if err != nil {
		lang = "en"
	}

	return editorRequestParams{token.User, token.Team, bid, fid, lang}, nil
}

func (c *editorController) fetchMiroData(ctx context.Context, params editorRequestParams, accessToken string) (*miro.BoardMemberResponse, *miro.FileInfoResponse, error) {
	g, _ := errgroup.WithContext(ctx)
	var userInfo *miro.BoardMemberResponse
	var fileInfo *miro.FileInfoResponse

	g.Go(func() error {
		var err error
		userInfo, err = c.MiroClient.GetBoardMember(ctx, miro.GetBoardMemberRequest{
			BoardID:  params.bid,
			MemberID: params.uid,
			Token:    accessToken,
		})

		if err != nil {
			return fmt.Errorf("failed to fetch user info: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		fileInfo, err = c.MiroClient.GetFileInfo(ctx, miro.GetFileInfoRequest{
			BoardID: params.bid,
			ItemID:  params.fid,
			Token:   accessToken,
		})

		if err != nil {
			return fmt.Errorf("failed to fetch file info: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	publicFile, err := c.MiroClient.GetFilePublicURL(ctx, miro.GetFilePublicURLRequest{
		URL:   fileInfo.Data.DocumentURL,
		Token: accessToken,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get public URL: %w", err)
	}

	fileInfo.Data.DocumentURL = publicFile.URL
	return userInfo, fileInfo, nil
}

func (c *editorController) resolveServerSettings(settings *component.Settings) (address, secret string, err error) {
	address = settings.Address
	secret = settings.Secret

	if settings.Demo.Enabled && address == "" && secret == "" {
		if settings.Demo.Started != nil {
			demoExpiry := settings.Demo.Started.Add(time.Duration(c.Config.DemoServer.Days) * 24 * time.Hour)
			if demoExpiry.After(time.Now()) {
				address = c.Config.DemoServer.Address
				secret = c.Config.DemoServer.Secret
			}
		}
	}

	if address == "" || secret == "" {
		return "", "", errors.New("please make sure you have configured the document server")
	}

	return address, secret, nil
}

func (c *editorController) buildEditorConfig(
	ctx context.Context,
	callbackURL string,
	boardID string,
	user *miro.BoardMemberResponse,
	file *miro.FileInfoResponse,
	secret string,
) (*document.Config, error) {
	config, err := c.BuilderService.Build(
		ctx,
		callbackURL,
		builderRequest{Board: boardID, File: *file},
		document.WithKey([]byte(secret)),
		document.WithUserConfigurer(user),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to build configuration: %w", err)
	}

	return config, nil
}

func (c *editorController) handleGet(ctx echo.Context) error {
	return c.ExecuteWithTimeout(ctx, 3*time.Second, func(tctx context.Context) error {
		handleRequestError := func(err error, message string) error {
			if err != nil {
				return c.HandleError(ctx, err, http.StatusBadRequest, message)
			}
			return nil
		}

		params, err := c.extractAndValidateParams(ctx)
		if err := handleRequestError(err, "failed to extract parameters"); err != nil {
			return err
		}

		settings, auth, err := c.FetchAuthenticationWithSettings(tctx, params.uid, params.tid, params.bid)
		if err := handleRequestError(err, "failed to fetch required data"); err != nil {
			return err
		}

		address, secret, err := c.resolveServerSettings(settings)
		if err := handleRequestError(err, "invalid configuration"); err != nil {
			return err
		}

		user, file, err := c.fetchMiroData(tctx, params, auth.AccessToken)
		if err := handleRequestError(err, "failed to fetch Miro data"); err != nil {
			return err
		}

		user.Lang = params.lang
		callbackURL := buildCallbackURL(c.Config.Server.CallbackURL, params.fid, params.uid, params.tid, params.bid)
		config, err := c.buildEditorConfig(tctx, callbackURL, params.bid, user, file, secret)
		if err := handleRequestError(err, "failed to build editor configuration"); err != nil {
			return err
		}

		configJSON, err := json.Marshal(config)
		if err := handleRequestError(err, "failed to encode configuration"); err != nil {
			return err
		}

		return ctx.Render(http.StatusOK, "editor", map[string]any{
			"apijs":  address + "/web-apps/apps/api/documents/api.js",
			"config": string(configJSON),
		})
	})
}

func (c *editorController) SupportedMethods() []common.HTTPMethod {
	return []common.HTTPMethod{common.MethodGet}
}
